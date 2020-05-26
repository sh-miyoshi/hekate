package oidc

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/user"
)

type option struct {
	audiences       []string
	genRefreshToken bool
	genIDToken      bool
	nonce           string
	maxAge          uint
}

// ClientAuth authenticates client with id and secret
func ClientAuth(projectName string, clientID string, clientSecret string) error {
	client, err := db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		e := errors.Cause(err)
		if e == model.ErrNoSuchClient || e == model.ErrClientValidateFailed {
			return errors.Wrap(ErrInvalidClient, err.Error())
		}
		return errors.Wrap(err, "Failed to get client")
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return errors.Wrap(ErrInvalidClient, "client auth failed")
		}
	}

	return nil
}

// ReqAuthByPassword ...
func ReqAuthByPassword(project *model.ProjectInfo, userName string, password string, r *http.Request) (*TokenResponse, error) {
	usr, err := user.Verify(project.Name, userName, password)
	if err != nil {
		if errors.Cause(err) == user.ErrAuthFailed {
			return nil, errors.Wrap(ErrRequestUnauthorized, "user authentication failed")
		}
		return nil, err
	}

	audiences := []string{usr.ID}
	clientID := r.Form.Get("client_id")
	if clientID != "" {
		audiences = append(audiences, clientID)
	}

	return genTokenRes(usr.ID, project, r, option{
		audiences:       audiences,
		genRefreshToken: true,
	})
}

// ReqAuthByCode ...
func ReqAuthByCode(project *model.ProjectInfo, clientID string, code string, r *http.Request) (*TokenResponse, error) {
	s, err := db.GetInst().AuthCodeSessionGetByCode(code)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchAuthCodeSession {
			// TODO(revoke all token in code.UserID) <- SHOULD
			return nil, errors.Wrap(ErrInvalidRequest, "no such code")
		}
		return nil, errors.Wrap(err, "Failed to get auth code info")
	}

	// Validate session info
	if time.Now().Unix() >= s.ExpiresIn.Unix() {
		return nil, errors.Wrap(ErrInvalidRequest, "code is already expired")
	}

	if s.ClientID != clientID {
		return nil, errors.Wrap(ErrRequestUnauthorized, "missing client id")
	}

	// Remove Authorized code
	if err := db.GetInst().AuthCodeSessionDelete(s.SessionID); err != nil {
		return nil, errors.Wrap(err, "Failed to delete auth code")
	}

	audiences := []string{
		s.UserID,
		s.ClientID,
	}

	return genTokenRes(s.UserID, project, r, option{
		audiences:       audiences,
		genRefreshToken: true,
		genIDToken:      true,
		nonce:           s.Nonce,
		maxAge:          s.MaxAge,
	})
}

// ReqAuthByRefreshToken ...
func ReqAuthByRefreshToken(project *model.ProjectInfo, clientID string, refreshToken string, r *http.Request) (*TokenResponse, error) {
	claims := &token.RefreshTokenClaims{}
	issuer := token.GetExpectIssuer(r)
	if err := token.ValidateRefreshToken(claims, refreshToken, issuer); err != nil {
		return nil, errors.Wrap(ErrRequestUnauthorized, fmt.Sprintf("Failed to verify token: %v", err))
	}

	ok := false
	for _, aud := range claims.Audience {
		if aud == clientID {
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.Wrap(ErrInvalidClient, "refresh token is not for the client")
	}

	// Delete previous token
	if err := db.GetInst().SessionDelete(claims.SessionID); err != nil {
		return nil, errors.Wrap(err, "Failed to revoke previous token")
	}

	return genTokenRes(claims.Subject, project, r, option{
		audiences:       claims.Audience,
		genRefreshToken: true,
	})
}

// ReqAuthByClientCredentials ...
func ReqAuthByClientCredentials(project *model.ProjectInfo, clientID string, r *http.Request) (*TokenResponse, error) {
	cli, err := db.GetInst().ClientGet(project.Name, clientID)
	if err != nil {
		return nil, errors.Wrap(err, "Get client info failed")
	}
	if cli.AccessType != "confidential" {
		return nil, ErrInvalidRequest
	}

	audiences := []string{
		clientID,
	}
	return genTokenRes("", project, r, option{
		audiences: audiences,
	})
}

func genTokenRes(userID string, project *model.ProjectInfo, r *http.Request, opt option) (*TokenResponse, error) {
	accessLifeSpan := project.TokenConfig.AccessTokenLifeSpan
	if opt.maxAge > 0 && opt.maxAge < accessLifeSpan {
		accessLifeSpan = opt.maxAge
	}

	// Generate JWT Token
	res := TokenResponse{
		TokenType: "Bearer",
		ExpiresIn: project.TokenConfig.AccessTokenLifeSpan,
	}

	accessTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(accessLifeSpan),
		ProjectName: project.Name,
		UserID:      userID,
	}

	audiences := []string{
		userID,
	}
	if len(opt.audiences) > 0 {
		audiences = opt.audiences
	}

	var err error
	res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate access token")
	}

	if opt.genRefreshToken {
		res.RefreshExpiresIn = project.TokenConfig.RefreshTokenLifeSpan
		if opt.maxAge > 0 {
			res.RefreshExpiresIn = opt.maxAge
		}

		refreshTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(res.RefreshExpiresIn),
			ProjectName: project.Name,
			UserID:      userID,
		}

		sessionID := uuid.New().String()
		res.RefreshToken, err = token.GenerateRefreshToken(sessionID, audiences, refreshTokenReq)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate refresh token")
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get IP")
		}
		ent := &model.Session{
			UserID:      userID,
			ProjectName: project.Name,
			SessionID:   sessionID,
			CreatedAt:   time.Now(),
			ExpiresIn:   res.RefreshExpiresIn,
			FromIP:      ip,
		}

		if err := db.GetInst().SessionAdd(ent); err != nil {
			return nil, errors.Wrap(err, "Failed to register session")
		}

	}

	if opt.genIDToken {
		lifeSpan := project.TokenConfig.AccessTokenLifeSpan
		var maxAge *uint
		if opt.maxAge > 0 {
			lifeSpan = opt.maxAge
			maxAge = &opt.maxAge
		}

		idTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(lifeSpan),
			ProjectName: project.Name,
			UserID:      userID,
			Nonce:       opt.nonce,
			MaxAge:      maxAge,
		}
		res.IDToken, err = token.GenerateIDToken(audiences, idTokenReq)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate id token")
		}
	}

	return &res, nil
}