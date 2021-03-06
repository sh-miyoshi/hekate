package authn

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/login"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
)

var (
	defaultScopes = []string{"openid", "email"}
)

type option struct {
	audiences       []string
	genRefreshToken bool
	genIDToken      bool
	nonce           string
	endUserAuthTime time.Time
	scopes          []string
}

// ReqAuthByPassword ...
func ReqAuthByPassword(project *model.ProjectInfo, userName string, password string, r *http.Request) (*oidc.TokenResponse, *errors.Error) {
	usr, err := login.UserVerifyByPassword(project.Name, userName, password)
	if err != nil {
		if errors.Contains(err, login.ErrAuthFailed) || errors.Contains(err, login.ErrUserLocked) {
			return nil, errors.Append(errors.ErrRequestUnauthorized, err.Error())
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
		endUserAuthTime: time.Unix(0, 0),
		scopes:          defaultScopes,
	})
}

// ReqAuthByCode ...
func ReqAuthByCode(project *model.ProjectInfo, clientID string, code string, codeVerifier string, r *http.Request) (*oidc.TokenResponse, *errors.Error) {
	s, err := db.GetInst().LoginSessionGetByCode(project.Name, code)
	if err != nil {
		// TODO(revoke all token in code.UserID) <- SHOULD
		if errors.Contains(err, model.ErrNoSuchLoginSession) {
			return nil, errors.Append(errors.ErrInvalidGrant, "no such code")
		}
		if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			return nil, errors.Append(errors.ErrInvalidRequest, "invalid login session code")
		}
		return nil, errors.Append(err, "Failed to get login session info")
	}

	// PKCE Code Verify
	if codeVerifier != "" {
		challenge := codeVerifier
		if s.CodeChallengeMethod == "S256" {
			// BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
			sum := sha256.Sum256([]byte(codeVerifier))
			challenge = base64.RawURLEncoding.EncodeToString(sum[:])
		}
		if challenge != s.CodeChallenge {
			logger.Debug("Expect PKCE code challenge: %s, but got %s", s.CodeChallenge, challenge)
			return nil, errors.Append(errors.ErrInvalidGrant, "code challenge failed")
		}
	} else {
		if s.CodeChallenge != "" {
			return nil, errors.Append(errors.ErrInvalidGrant, "code challenge is not registed")
		}
	}

	// Validate session info
	if time.Now().After(s.ExpiresDate) {
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
		endUserAuthTime: s.LoginDate,
		scopes:          s.Scopes,
	})
}

// ReqAuthByRefreshToken ...
func ReqAuthByRefreshToken(project *model.ProjectInfo, clientID string, refreshToken string, r *http.Request) (*oidc.TokenResponse, *errors.Error) {
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
		endUserAuthTime: s.LastAuthTime,
		scopes:          strings.Split(claims.Scope, " "),
	})
}

// ReqAuthByClientCredentials ...
func ReqAuthByClientCredentials(project *model.ProjectInfo, clientID string, r *http.Request) (*oidc.TokenResponse, *errors.Error) {
	cli, err := db.GetInst().ClientGet(project.Name, clientID)
	if err != nil {
		return nil, errors.Append(err, "Get client info failed")
	}
	if cli.AccessType != "confidential" {
		return nil, errors.Append(errors.ErrInvalidRequest, "client type is not confidential")
	}

	audiences := []string{
		clientID,
	}
	return genTokenRes("", project, r, option{
		audiences: audiences,
	})
}

// ReqAuthByDeviceCode ...
func ReqAuthByDeviceCode(project *model.ProjectInfo, clientID string, deviceCode string, r *http.Request) (*oidc.TokenResponse, *errors.Error) {
	devices, err := db.GetInst().DeviceGetList(project.Name, &model.DeviceFilter{DeviceCode: deviceCode})
	if err != nil {
		return nil, errors.Append(err, "Get device failed")
	}
	if len(devices) == 0 {
		return nil, errors.ErrInvalidRequest
	}

	device := devices[0]

	// get login session
	logger.Debug("login session ID of device authentication: %s", device.LoginSessionID)
	s, err := db.GetInst().LoginSessionGet(project.Name, device.LoginSessionID)
	if err != nil {
		return nil, errors.Append(err, "Failed to get login session info")
	}
	if s.LoginDate.IsZero() {
		logger.Debug("user is not logged in yet")
		return nil, errors.ErrAuthorizationPending
	}

	// user already logged in, so delete login session and return token
	if err := db.GetInst().DeviceDelete(project.Name, device.DeviceCode); err != nil {
		return nil, errors.Append(err, "Failed to delete device")
	}

	if time.Now().After(s.ExpiresDate) {
		logger.Info("the device session was already expired")
		return nil, errors.ErrExpiredToken
	}

	audiences := []string{
		clientID,
	}
	return genTokenRes(s.UserID, project, r, option{
		audiences:       audiences,
		genRefreshToken: true,
		endUserAuthTime: s.LoginDate,
		scopes:          s.Scopes,
	})
}

func genTokenRes(userID string, project *model.ProjectInfo, r *http.Request, opt option) (*oidc.TokenResponse, *errors.Error) {
	// Generate JWT Token
	res := oidc.TokenResponse{
		TokenType: "Bearer",
		ExpiresIn: project.TokenConfig.AccessTokenLifeSpan,
	}

	accessTokenReq := token.Request{
		Issuer:      token.GetFullIssuer(r),
		ExpiresIn:   int64(project.TokenConfig.AccessTokenLifeSpan),
		ProjectName: project.Name,
		UserID:      userID,
		Scopes:      opt.scopes,
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
		refreshTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiresIn:   int64(res.RefreshExpiresIn),
			ProjectName: project.Name,
			UserID:      userID,
			Scopes:      opt.scopes,
		}

		sessionID := uuid.New().String()
		res.RefreshToken, err = token.GenerateRefreshToken(sessionID, audiences, refreshTokenReq)
		if err != nil {
			return nil, errors.Append(err, "Failed to generate refresh token")
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, errors.New("Invalid request", "Failed to get IP: %v", err)
		}
		ent := &model.Session{
			UserID:       userID,
			ProjectName:  project.Name,
			SessionID:    sessionID,
			CreatedAt:    time.Now(),
			ExpiresIn:    int64(res.RefreshExpiresIn),
			FromIP:       ip,
			LastAuthTime: opt.endUserAuthTime,
		}

		if err := db.GetInst().SessionAdd(project.Name, ent); err != nil {
			return nil, errors.Append(err, "Failed to register session")
		}

	}

	if opt.genIDToken {
		idTokenReq := token.Request{
			Issuer:          token.GetFullIssuer(r),
			ExpiresIn:       int64(project.TokenConfig.AccessTokenLifeSpan),
			ProjectName:     project.Name,
			UserID:          userID,
			Nonce:           opt.nonce,
			EndUserAuthTime: opt.endUserAuthTime,
			Scopes:          opt.scopes,
		}
		res.IDToken, err = token.GenerateIDToken(audiences, idTokenReq)
		if err != nil {
			return nil, errors.Append(err, "Failed to generate id token")
		}
	}

	return &res, nil
}
