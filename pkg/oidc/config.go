package oidc

import (
	"path/filepath"
	"time"
)

var (
	expiresTimeSec     uint64
	userLoginHTML      string
	userLoginErrorHTML string
	userConsentHTML    string
	userLoginResPath   string
	cookieSecure       bool
)

// InitConfig ...
func InitConfig(runAsHTTPS bool, loginSessionExpiresTimeSec uint64, loginResDir string, resourcePath string) {
	userLoginHTML = filepath.Join(loginResDir, "index.html")
	userLoginErrorHTML = filepath.Join(loginResDir, "error.html")
	userConsentHTML = filepath.Join(loginResDir, "consent.html")
	expiresTimeSec = loginSessionExpiresTimeSec
	userLoginResPath = resourcePath
	cookieSecure = runAsHTTPS
}

// GetSupportedResponseType ...
func GetSupportedResponseType() []string {
	return []string{
		"code",
		"id_token",
		"token",
		"code id_token",
		"code token",
		"id_token token",
		"code id_token token",
		// TODO(support type "none")
	}
}

// GetSupportedScope ...
func GetSupportedScope() []string {
	return []string{"openid"}
}

// GetLoginSessionExpiresTime ...
func GetLoginSessionExpiresTime() time.Duration {
	return time.Second * time.Duration(expiresTimeSec)
}

// IsCookieSecure ...
func IsCookieSecure() bool {
	return cookieSecure
}
