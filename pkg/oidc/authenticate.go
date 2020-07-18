package oidc

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/user"
)

type option struct {
	audiences       []string
	genRefreshToken bool
	genIDToken      bool
	nonce           string
	maxAge          uint
	endUserAuthTime time.Time
}

// ClientAuth authenticates client with id and secret
func ClientAuth(projectName string, clientID string, clientSecret string) *errors.Error {
	client, err := db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchClient) || errors.Contains(err, model.ErrClientValidateFailed) {
			return errors.Append(errors.ErrInvalidClient, err.Error())
		}
		return errors.Append(err, "Failed to get client")
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return errors.Append(errors.ErrInvalidClient, "client auth failed")
		}
	}

	return nil
}

// ReqAuthByPassword ...
func ReqAuthByPassword(project *model.ProjectInfo, userName string, password string, r *http.Request) (*TokenResponse, *errors.Error) {
	usr, err := user.Verify(project.Name, userName, password)
	if err != nil {
		if errors.Contains(err, user.ErrAuthFailed) {
			return nil, errors.Append(errors.ErrRequestUnauthorized, "user authentication failed")
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
		maxAge:          0,
		endUserAuthTime: time.Unix(0, 0),
	})
}

// ReqAuthByCode ...
func ReqAuthByCode(project *model.ProjectInfo, clientID string, code string, r *http.Request) (*TokenResponse, *errors.Error) {
	s, err := db.GetInst().LoginSessionGetByCode(project.Name, code)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchLoginSession) {
			// TODO(revoke all token in code.UserID) <- SHOULD
			return nil, errors.Append(errors.ErrInvalidGrant, "no such code")
		}
		return nil, errors.Append(err, "Failed to get login session info")
	}

	// Validate session info
	if time.Now().Unix() >= s.ExpiresIn.Unix() {
		return nil, errors.Append(errors.ErrInvalidRequest, "code is already expired")
	}

	if s.ClientID != clientID {
		return nil, errors.Append(errors.ErrRequestUnauthorized, "missing client id")
	}

	// Remove Authorized code
	if err := db.GetInst().LoginSessionDelete(project.Name, s.SessionID); err != nil {
		return nil, errors.Append(err, "Failed to delete login session")
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
		endUserAuthTime: s.LoginDate,
	})
}

// ReqAuthByRefreshToken ...
func ReqAuthByRefreshToken(project *model.ProjectInfo, clientID string, refreshToken string, r *http.Request) (*TokenResponse, *errors.Error) {
	claims := &token.RefreshTokenClaims{}
	issuer := token.GetExpectIssuer(r)
	if err := token.ValidateRefreshToken(claims, refreshToken, issuer); err != nil {
		return nil, errors.Append(errors.ErrRequestUnauthorized, fmt.Sprintf("Failed to verify token: %v", err))
	}

	ok := false
	for _, aud := range claims.Audience {
		if aud == clientID {
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.Append(errors.ErrInvalidClient, "refresh token is not for the client")
	}

	s, err := db.GetInst().SessionGet(project.Name, claims.SessionID)
	if err != nil {
		return nil, errors.Append(err, "Failed to get previous token")
	}

	// Delete previous token
	if err := db.GetInst().SessionDelete(project.Name, claims.SessionID); err != nil {
		return nil, errors.Append(err, "Failed to revoke previous token")
	}

	return genTokenRes(claims.Subject, project, r, option{
		audiences:       claims.Audience,
		genRefreshToken: true,
		maxAge:          s.AuthMaxAge,
		endUserAuthTime: s.LastAuthTime,
	})
}

// ReqAuthByClientCredentials ...
func ReqAuthByClientCredentials(project *model.ProjectInfo, clientID string, r *http.Request) (*TokenResponse, *errors.Error) {
	cli, err := db.GetInst().ClientGet(project.Name, clientID)
	if err != nil {
		return nil, errors.Append(err, "Get client info failed")
	}
	if cli.AccessType != "confidential" {
		return nil, errors.ErrInvalidRequest
	}

	audiences := []string{
		clientID,
	}
	return genTokenRes("", project, r, option{
		audiences: audiences,
	})
}

func genTokenRes(userID string, project *model.ProjectInfo, r *http.Request, opt option) (*TokenResponse, *errors.Error) {
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

	var err *errors.Error
	res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
	if err != nil {
		return nil, errors.Append(err, "Failed to generate access token")
	}

	if opt.genRefreshToken {
		res.RefreshExpiresIn = project.TokenConfig.RefreshTokenLifeSpan
		maxAge := project.TokenConfig.AccessTokenLifeSpan
		if opt.maxAge > 0 {
			res.RefreshExpiresIn = opt.maxAge
			maxAge = opt.maxAge
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
			return nil, errors.Append(err, "Failed to generate refresh token")
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, errors.New("", "Failed to get IP: %v", err)
		}
		ent := &model.Session{
			UserID:       userID,
			ProjectName:  project.Name,
			SessionID:    sessionID,
			CreatedAt:    time.Now(),
			ExpiresIn:    res.RefreshExpiresIn,
			FromIP:       ip,
			LastAuthTime: opt.endUserAuthTime,
			AuthMaxAge:   maxAge,
		}

		if err := db.GetInst().SessionAdd(project.Name, ent); err != nil {
			return nil, errors.Append(err, "Failed to register session")
		}

	}

	if opt.genIDToken {
		lifeSpan := project.TokenConfig.AccessTokenLifeSpan
		if opt.maxAge > 0 {
			lifeSpan = opt.maxAge
		}

		idTokenReq := token.Request{
			Issuer:          token.GetFullIssuer(r),
			ExpiredTime:     time.Second * time.Duration(lifeSpan),
			ProjectName:     project.Name,
			UserID:          userID,
			Nonce:           opt.nonce,
			EndUserAuthTime: opt.endUserAuthTime,
		}
		res.IDToken, err = token.GenerateIDToken(audiences, idTokenReq)
		if err != nil {
			return nil, errors.Append(err, "Failed to generate id token")
		}
	}

	return &res, nil
}
