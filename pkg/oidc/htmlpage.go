package oidc

import (
	"html/template"
	"net/http"

	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// WriteUserLoginPage ...
func WriteUserLoginPage(projectName, sessionID, errMsg, state string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(userLoginHTML)
	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		http.Error(w, "User Login Page maybe broken", http.StatusInternalServerError)
		return
	}

	url := "/api/v1/project/" + projectName + "/openid-connect/login?login_session_id=" + sessionID
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
