package oidc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net"
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

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get from IP: %v", err)
	}
	ent := &model.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresIn: res.RefreshExpiresIn,
		FromIP:    ip,
	}

	if err := db.GetInst().SessionAdd(ent); err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to register session: %+v", err)
	}

	return &res, http.StatusOK, ""
}

// AuthByCode ...
func AuthByCode(project *model.ProjectInfo, codeID string, r *http.Request) (*TokenResponse, int, string) {
	code, status, msg := ValidateAuthCode(codeID)
	if status != http.StatusOK {
		return nil, status, msg
	}

	if err := db.GetInst().AuthCodeDelete(codeID); err != nil {
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

// AuthByRefreshToken ...
func AuthByRefreshToken(project *model.ProjectInfo, refreshToken string, r *http.Request) (*TokenResponse, int, string) {
	clientID := r.Form.Get("client_id")

	claims := &token.RefreshTokenClaims{}
	issuer := token.GetExpectIssuer(r)
	if err := token.ValidateRefreshToken(claims, refreshToken, issuer); err != nil {
		return nil, http.StatusBadRequest, fmt.Sprintf("Failed to validate token: %v", err)
	}

	ok := false
	for _, aud := range claims.Audience {
		if aud == clientID {
			ok = true
			break
		}
	}

	if !ok {
		return nil, http.StatusBadRequest, "refresh token is not for the client"
	}

	// Revoke previous token
	if err := db.GetInst().SessionDelete(claims.SessionID); err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to revoke previous token: %+v", err)
	}

	// Generate JWT Token
	res := TokenResponse{
		TokenType:        "Bearer",
		ExpiresIn:        project.TokenConfig.AccessTokenLifeSpan,
		RefreshExpiresIn: project.TokenConfig.RefreshTokenLifeSpan,
	}

	userID := claims.Subject

	accessTokenReq := token.Request{
		Issuer:      claims.Issuer,
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: claims.Project,
		UserID:      userID,
	}

	var err error
	res.AccessToken, err = token.GenerateAccessToken(claims.Audience, accessTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	refreshTokenReq := token.Request{
		Issuer:      claims.Issuer,
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
		ProjectName: claims.Project,
		UserID:      userID,
	}

	sessionID := uuid.New().String()
	res.RefreshToken, err = token.GenerateRefreshToken(sessionID, claims.Audience, refreshTokenReq)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get JWT token: %+v", err)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to get from IP: %v", err)
	}
	ent := &model.Session{
		UserID:    userID,
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresIn: res.RefreshExpiresIn,
		FromIP:    ip,
	}
	if err := db.GetInst().SessionAdd(ent); err != nil {
		return nil, http.StatusInternalServerError, fmt.Sprintf("Failed to register refresh token session token: %+v", err)
	}

	return &res, http.StatusOK, ""
}
