package config

import (
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
	"os"
)

type secret struct {
	AccessToken string `json:"accessToken"`
	AccessTokenExpiresTime time.Time `json:"accessTokenExpiresTime"`
	RefreshToken string `json:"refreshToken"`
	RefreshTokenExpiresTime time.Time `json:"refreshTokenExpiresTime"`
}

// SetSecret ...
func SetSecret(token *tokenapi.TokenResponse) {
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")

	v := secret{
		AccessToken: token.AccessToken,
		AccessTokenExpiresTime: time.Now().Add(time.Second*time.Duration(token.AccessExpiresIn)),
		RefreshToken: token.RefreshToken,
		RefreshTokenExpiresTime: time.Now().Add(time.Second*time.Duration(token.RefreshExpiresIn)),
	}

	bytes, _ := json.MarshalIndent(v, "", "  ")
	ioutil.WriteFile(secretFile, bytes, os.ModePerm)
}

// GetAccessToken ...
func GetAccessToken() (string, error) {
	// Get Secret Info
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")
	buf, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read secret file: %v\nYou need to `jwtctl login` at first", err)
	}

	var s secret
	json.Unmarshal(buf, &s)

	// TODO(Validate secret)
	// TODO(Refresh token if required)

	return s.AccessToken, nil
}