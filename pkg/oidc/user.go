package oidc

import (
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

type sessionInfo struct {
	VerifyCode  string
	ExpiresIn   time.Time
	BaseRequest *AuthRequest
}

var (
	// TODO(use database for scale)
	loginSessions = make(map[string]*sessionInfo) // key: verifyCode
)

// WriteUserLoginPage ...
func WriteUserLoginPage(code string, projectName string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		http.Error(w, "User Login Page maybe broken", http.StatusInternalServerError)
		return
	}

	url := "/api/v1/project/" + projectName + "/openid-connect/login?login_verify_code=" + code

	d := map[string]string{
		"URL": url,
	}

	tpl.Execute(w, d)
}

// RegisterUserLoginSession ...
func RegisterUserLoginSession(req *AuthRequest) string {
	code := uuid.New().String()
	loginSessions[code] = &sessionInfo{
		VerifyCode:  code,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		BaseRequest: req,
	}
	return code
}

// UserLoginVerify ...
func UserLoginVerify(code string) (*AuthRequest, error) {
	if s, ok := loginSessions[code]; ok {
		// Get is only allowed once
		delete(loginSessions, code)

		now := time.Now().Unix()
		if now > s.ExpiresIn.Unix() {
			return nil, errors.New("Session already expired")
		}
		return s.BaseRequest, nil
	}
	return nil, errors.New("No such session")
}
