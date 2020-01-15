package oidc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"time"
)

// TokenResponse ...
type TokenResponse struct {
	TokenType        string
	AccessToken      string
	ExpiresIn        uint
	RefreshToken     string
	RefreshExpiresIn uint
	IDToken          string
}

// AuthByPassword ...
func AuthByPassword(project *model.ProjectInfo, userName string, password string, r *http.Request) (*TokenResponse, int, string) {
	user, err := db.GetInst().UserGetByName(project.Name, userName)
	if err != nil {
		if err == model.ErrNoSuchUser {
			return nil, http.StatusBadRequest, fmt.Sprintf("No such user %s in project %s", userName, project.Name)
		}
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get user id: %+v", err)
	}

	hash := util.CreateHash(password)
	if user.PasswordHash != hash {
		return nil, http.StatusBadRequest, "password authentication failed"
	}

	// Generate JWT Token
	res := TokenResponse{
		TokenType:        "Bearer",
		ExpiresIn:        project.TokenConfig.AccessTokenLifeSpan,
		RefreshExpiresIn: project.TokenConfig.RefreshTokenLifeSpan,
	}

	audiences := []string{user.ID}
	clientID := r.Form.Get("client_id")
	if clientID != "" {
		audiences = append(audiences, clientID)
	}

	accessTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      user.ID,
	}

	res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	refreshTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      user.ID,
	}

	sessionID := uuid.New().String()
	res.RefreshToken, err = token.GenerateRefreshToken(sessionID, audiences, refreshTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	if err := db.GetInst().NewSession(user.ID, sessionID, res.RefreshExpiresIn, r.RemoteAddr); err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to register refresh token session token: %+v", err)
	}

	return &res, http.StatusOK, ""
}

// AuthByCode ...
func AuthByCode(project *model.ProjectInfo, codeID string, r *http.Request) (*TokenResponse, int, string) {
	code, status, msg := ValidateAuthCode(codeID)
	if status != http.StatusOK {
		return nil, status, msg
	}

	if err := db.GetInst().DeleteAuthCode(codeID); err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to delete auth code: %+v", err)
	}

	res := TokenResponse{
		TokenType: "Bearer",
		ExpiresIn: project.TokenConfig.AccessTokenLifeSpan,
	}

	audiences := []string{
		code.UserID,
		code.ClientID,
	}

	accessTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      code.UserID,
	}

	var err error
	res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	audiences = []string{
		code.ClientID,
	}

	idTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      code.UserID,
	}
	res.IDToken, err = token.GenerateRefreshToken("", audiences, idTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	return &res, http.StatusOK, ""
}
