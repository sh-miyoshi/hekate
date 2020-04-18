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

func genTokenRes(audiences []string, userID string, project *model.ProjectInfo, r *http.Request, genRefresh, genIDToken bool, nonce string) (*TokenResponse, error) {
	// Generate JWT Token
	res := TokenResponse{
		TokenType: "Bearer",
		ExpiresIn: project.TokenConfig.AccessTokenLifeSpan,
	}

	accessTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      userID,
	}

	var err error
	res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate access token")
	}

	if genRefresh {
		res.RefreshExpiresIn = project.TokenConfig.RefreshTokenLifeSpan
		refreshTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
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

	if genIDToken {
		idTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
			ProjectName: project.Name,
			UserID:      userID,
			Nonce:       nonce,
		}
		res.IDToken, err = token.GenerateIDToken("", audiences, idTokenReq)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate id token")
		}
	}

	return &res, nil
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

	return genTokenRes(audiences, usr.ID, project, r, true, false, "")
}

// ReqAuthByCode ...
func ReqAuthByCode(project *model.ProjectInfo, clientID string, codeID string, r *http.Request) (*TokenResponse, error) {
	code, err := verifyAuthCode(codeID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify code")
	}

	if code.ClientID != clientID {
		return nil, errors.Wrap(ErrRequestUnauthorized, "missing client id")
	}

	// Remove Authorized code
	if err := db.GetInst().AuthCodeDelete(codeID); err != nil {
		return nil, errors.Wrap(err, "Failed to delete auth code")
	}

	audiences := []string{
		code.UserID,
		code.ClientID,
	}

	return genTokenRes(audiences, code.UserID, project, r, true, true, code.Nonce)
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

	userID := claims.Subject
	return genTokenRes(claims.Audience, userID, project, r, true, false, "")
}

// ReqAuthByRClientCredentials ...
func ReqAuthByRClientCredentials(project *model.ProjectInfo, clientID string, r *http.Request) (*TokenResponse, error) {
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
	return genTokenRes(audiences, "", project, r, false, false, "")
}
