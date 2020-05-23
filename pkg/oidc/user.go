package oidc

import (
	"html/template"
	"net/http"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// WriteUserLoginPage ...
func WriteUserLoginPage(projectName, code, errMsg, state string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		http.Error(w, "User Login Page maybe broken", http.StatusInternalServerError)
		return
	}

	url := "/api/v1/project/" + projectName + "/openid-connect/login?login_verify_code=" + code
	if state != "" {
		url += "&state=" + state
	}

	d := map[string]string{
		"URL":             url,
		"CSSResourcePath": userLoginResPath + "/css",
		"IMGResourcePath": userLoginResPath + "/img",
		"Error":           errMsg,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

// WriteErrorPage ...
func WriteErrorPage(errMsg string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginErrorHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		http.Error(w, "User Login Error Page maybe broken", http.StatusInternalServerError)
		return
	}

	d := map[string]string{
		"CSSResourcePath": userLoginResPath + "/css",
		"IMGResourcePath": userLoginResPath + "/img",
		"Error":           errMsg,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

// WriteConsentPage ...
func WriteConsentPage(w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userConsentHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		http.Error(w, "User Consent Page maybe broken", http.StatusInternalServerError)
		return
	}

	// TODO(set this)
	url := ""
	d := map[string]string{
		"CSSResourcePath": userLoginResPath + "/css",
		"IMGResourcePath": userLoginResPath + "/img",
		"URL":             url,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

// RegisterUserLoginSession ...
func RegisterUserLoginSession(projectName string, req *AuthRequest) (string, error) {
	code := uuid.New().String()

	resMode := req.ResponseMode
	if resMode == "" {
		resMode = "fragment"
		if len(req.ResponseType) == 1 {
			if req.ResponseType[0] == "code" || req.ResponseType[0] == "none" {
				resMode = "query"
			}
		}
	}

	s := &model.LoginSessionInfo{
		VerifyCode:   code,
		ExpiresIn:    time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		ClientID:     req.ClientID,
		RedirectURI:  req.RedirectURI,
		Nonce:        req.Nonce,
		ProjectName:  projectName,
		MaxAge:       req.MaxAge,
		ResponseMode: resMode,
		ResponseType: req.ResponseType,
		Prompt:       req.Prompt,
	}

	if err := db.GetInst().LoginSessionAdd(s); err != nil {
		return "", errors.Wrap(err, "add user login session failed")
	}
	return code, nil
}

// UserLoginVerify ...
func UserLoginVerify(code string) (*UserLoginInfo, error) {
	s, err := db.GetInst().LoginSessionGet(code)
	if err != nil {
		return nil, errors.Wrap(err, "user login session get failed")
	}

	if err := db.GetInst().LoginSessionDelete(code); err != nil {
		return nil, errors.Wrap(err, "user login sessiond delete failed")
	}

	now := time.Now().Unix()
	if now > s.ExpiresIn.Unix() {
		return nil, errors.New("Session already expired")
	}
	return &UserLoginInfo{
		Scope:        s.Scope,
		ResponseType: s.ResponseType,
		ClientID:     s.ClientID,
		RedirectURI:  s.RedirectURI,
		Nonce:        s.Nonce,
		MaxAge:       s.MaxAge,
		ResponseMode: s.ResponseMode,
		Prompt:       s.Prompt,
	}, nil
}
