package sso

import (
	"crypto/x509"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
)

// SetSSOSessionToCookie ...
func SetSSOSessionToCookie(w http.ResponseWriter, projectName, userID, issuer string) *errors.Error {
	cfg := config.Get()

	req := token.Request{
		Issuer:      issuer,
		ExpiredTime: time.Second * time.Duration(cfg.SSOExpiresTime),
		ProjectName: projectName,
		UserID:      userID,
	}
	tkn, err := token.GenerateSSOToken(req)
	if err != nil {
		return errors.Append(err, "Failed to generate SSO token")
	}

	cookie := &http.Cookie{
		Name:     "HEKATE_LOGIN_SESSION",
		Value:    tkn,
		MaxAge:   int(req.ExpiredTime),
		Secure:   cfg.HTTPSConfig.Enabled,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return nil
}

// GetLoginUserIDFromSSOSessionCookie ...
func GetLoginUserIDFromSSOSessionCookie(cookie *http.Cookie, projectName string) (string, *errors.Error) {
	var claims jwt.StandardClaims
	tkn, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(projectName)
		if err != nil {
			return nil, errors.Append(err, "Failed to get project")
		}

		switch token.Method {
		case jwt.SigningMethodRS256:
			key, err := x509.ParsePKCS1PublicKey(project.TokenConfig.SignPublicKey)
			if err != nil {
				return nil, errors.New("Invalid request", "Failed to parse public key: %v", err)
			}
			return key, nil
		}

		return nil, errors.New("Invalid request", "unknown token sigining method")
	})

	if err != nil || !tkn.Valid {
		return "", errors.New("Invalid request", "Token in cookie is not valid: %v", err)
	}

	return claims.Subject, nil
}

// Handle method return redirect page after logged in when found valid session
func Handle(method string, projectName string, userID string, tokenIssuer string, authReq *oidc.AuthRequest) (*http.Request, *errors.Error) {
	sessions, err := db.GetInst().SessionGetList(projectName, &model.SessionFilter{UserID: userID})
	if err != nil {
		return nil, errors.Append(err, "Failed to get session list")
	}
	if len(sessions) == 0 {
		return nil, errors.Append(errors.ErrLoginRequired, "No sessions, so return login_required")
	}

	// check max_age
	// if now > auth_time + max_age return login_required
	now := time.Now()
	for _, s := range sessions {
		lifeSpan := s.ExpiresIn
		if authReq.MaxAge > 0 {
			lifeSpan = authReq.MaxAge
		}

		valid := s.LastAuthTime.Add(time.Second * time.Duration(lifeSpan))
		logger.Debug("Session Info: now %v, valid time %v", now, valid)
		if now.Before(valid) {
			expires := time.Second * time.Duration(config.Get().LoginSessionExpiresTime)

			ls := &model.LoginSession{
				SessionID:           uuid.New().String(),
				ResponseType:        authReq.ResponseType,
				ProjectName:         projectName,
				UserID:              s.UserID,
				ClientID:            authReq.ClientID,
				Nonce:               authReq.Nonce,
				LoginDate:           s.LastAuthTime,
				RedirectURI:         authReq.RedirectURI,
				ResponseMode:        authReq.ResponseMode,
				Scope:               authReq.Scope,
				Prompt:              authReq.Prompt,
				ExpiresIn:           time.Now().Add(expires).Unix(),
				CodeChallenge:       authReq.CodeChallenge,
				CodeChallengeMethod: authReq.CodeChallengeMethod,
			}
			req, err := oidc.CreateLoggedInResponse(ls, authReq.State, tokenIssuer)
			if err != nil {
				return nil, errors.Append(err, "Failed to create login redirect info")
			}
			if err := db.GetInst().LoginSessionAdd(projectName, ls); err != nil {
				return nil, errors.Append(err, "Failed to register login session")
			}
			return req, nil
		}
	}

	return nil, errors.Append(errors.ErrLoginRequired, "No valid session, so return login_required")
}
