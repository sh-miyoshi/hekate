package oidc

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/user"
	"net"
	"net/http"
	"time"
)

type sessionInfo struct {
	VerifyCode  string
	ExpiresIn   time.Time
	BaseRequest *AuthRequest
}

var (
	expiresTimeSec uint64
	userLoginHTML  string
	loginSessions  map[string]*sessionInfo // key: verifyCode

	// ErrAuthCodeVerifyFailed ...
	ErrAuthCodeVerifyFailed = errors.New("failed to verify code")
)

func genTokenRes(audiences []string, userID string, project *model.ProjectInfo, r *http.Request, genRefresh, genIDToken bool) (*TokenResponse, error) {
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
			UserID:    userID,
			SessionID: sessionID,
			CreatedAt: time.Now(),
			ExpiresIn: res.RefreshExpiresIn,
			FromIP:    ip,
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
		}
		res.IDToken, err = token.GenerateRefreshToken("", audiences, idTokenReq)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate id token")
		}
	}

	return &res, nil
}

func verifyAuthCode(codeID string) (*model.AuthCode, error) {
	code, err := db.GetInst().AuthCodeGet(codeID)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchCode {
			// TODO(revoke all token in code.UserID) <- SHOULD
			return nil, errors.Wrap(ErrAuthCodeVerifyFailed, "no such code")
		}
		return nil, err
	}
	logger.Debug("Code: %v", code)

	if time.Now().Unix() >= code.ExpiresIn.Unix() {
		return nil, errors.Wrap(ErrAuthCodeVerifyFailed, "code is already expired")
	}

	return code, nil
}

// InitAuthCodeConfig ...
func InitAuthCodeConfig(authCodeExpiresTimeSec uint64, authCodeUserLoginFile string) {
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginHTML = authCodeUserLoginFile
	loginSessions = make(map[string]*sessionInfo)
}

// ReqAuthByPassword ...
func ReqAuthByPassword(project *model.ProjectInfo, userName string, password string, r *http.Request) (*TokenResponse, error) {
	usr, err := user.Verify(project.Name, userName, password)
	if err != nil {
		return nil, err
	}

	audiences := []string{usr.ID}
	clientID := r.Form.Get("client_id")
	if clientID != "" {
		audiences = append(audiences, clientID)
	}

	return genTokenRes(audiences, usr.ID, project, r, true, false)
}

// ReqAuthByCode ...
func ReqAuthByCode(project *model.ProjectInfo, codeID string, r *http.Request) (*TokenResponse, error) {
	code, err := verifyAuthCode(codeID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify code")
	}

	// Remove Authorized code
	if err := db.GetInst().AuthCodeDelete(codeID); err != nil {
		return nil, errors.Wrap(err, "Failed to delete auth code")
	}

	audiences := []string{
		code.UserID,
		code.ClientID,
	}

	return genTokenRes(audiences, code.UserID, project, r, false, true)
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

// GenerateAuthCode ...
func GenerateAuthCode(clientID string, redirectURL string, userID string) (string, error) {
	code := &model.AuthCode{
		CodeID:      uuid.New().String(),
		ClientID:    clientID,
		RedirectURL: redirectURL,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		UserID:      userID,
	}

	err := db.GetInst().AuthCodeAdd(code)

	return code.CodeID, err
}
