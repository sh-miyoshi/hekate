package oidc

import (
	"html/template"
	"net/http"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// WriteUserLoginPage ...
func WriteUserLoginPage(projectName, sessionID, errMsg, state string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		errors.WriteHTTPError(w, "Page Broken", errors.New("User Login Page maybe broken", ""), http.StatusInternalServerError)
		return
	}

	url := "/api/v1/project/" + projectName + "/openid-connect/login?login_session_id=" + sessionID
	if state != "" {
		url += "&state=" + state
	}

	d := map[string]string{
		"URL":                url,
		"StaticResourcePath": userLoginResPath + "/static",
		"Error":              errMsg,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

// WriteErrorPage ...
func WriteErrorPage(errMsg string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginErrorHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		errors.WriteHTTPError(w, "Page Broken", errors.New("User Login Error Page maybe broken", ""), http.StatusInternalServerError)
		return
	}

	d := map[string]string{
		"StaticResourcePath": userLoginResPath + "/static",
		"Error":              errMsg,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

// WriteConsentPage ...
func WriteConsentPage(projectName, sessionID, state string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userConsentHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		errors.WriteHTTPError(w, "Page Broken", errors.New("User Login Consent Page maybe broken", ""), http.StatusInternalServerError)
		return
	}

	url := "/api/v1/project/" + projectName + "/openid-connect/consent?login_session_id=" + sessionID
	if state != "" {
		url += "&state=" + state
	}

	d := map[string]string{
		"StaticResourcePath": userLoginResPath + "/static",
		"URL":                url,
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}
