package oidc

import (
	"path/filepath"
)

var (
	expiresTimeSec     uint64
	userLoginHTML      string
	userLoginErrorHTML string
	userLoginResPath   string
)

// InitConfig ...
func InitConfig(authCodeExpiresTimeSec uint64, loginResDir string, resourcePath string) {
	userLoginHTML = filepath.Join(loginResDir, "index.html")
	userLoginErrorHTML = filepath.Join(loginResDir, "error.html")
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginResPath = resourcePath
}

// GetSupportedResponseType ...
func GetSupportedResponseType() []string {
	return []string{
		"code",
		"id_token",
		// TODO(must be supported 'token id_token')
	}
}
