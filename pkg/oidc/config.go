package oidc

import (
	"path/filepath"
)

var (
	expiresTimeSec     uint64
	userLoginHTML      string
	userLoginErrorHTML string
	userConsentHTML    string
	userLoginResPath   string
)

// InitConfig ...
func InitConfig(authCodeExpiresTimeSec uint64, loginResDir string, resourcePath string) {
	userLoginHTML = filepath.Join(loginResDir, "index.html")
	userLoginErrorHTML = filepath.Join(loginResDir, "error.html")
	userConsentHTML = filepath.Join(loginResDir, "consent.html")
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginResPath = resourcePath
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
