package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
	"github.com/sh-miyoshi/hekate/pkg/hctl/login"
)

type secret struct {
	UserName                string    `json:"userName"`
	AccessToken             string    `json:"accessToken"`
	AccessTokenExpiresTime  time.Time `json:"accessTokenExpiresTime"`
	RefreshToken            string    `json:"refreshToken"`
	RefreshTokenExpiresTime time.Time `json:"refreshTokenExpiresTime"`
}

// SetSecret ...
func SetSecret(userName string, token *oidcapi.TokenResponse) {
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")

	v := secret{
		UserName:                userName,
		AccessToken:             token.AccessToken,
		AccessTokenExpiresTime:  time.Now().Add(time.Second * time.Duration(token.ExpiresIn)),
		RefreshToken:            token.RefreshToken,
		RefreshTokenExpiresTime: time.Now().Add(time.Second * time.Duration(token.RefreshExpiresIn)),
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
		return "", fmt.Errorf("Failed to read secret file: %v\nYou need to `hctl login` at first", err)
	}

	var s secret
	json.Unmarshal(buf, &s)

	if time.Now().After(s.RefreshTokenExpiresTime) {
		return "", fmt.Errorf("Token is expired\nPlease run `hctl login`")
	}

	if time.Now().After(s.AccessTokenExpiresTime) {
		// Refresh token by using refresh-token
		res, err := login.DoWithRefresh(sysConf.ServerAddr, sysConf.ProjectName, s.RefreshToken)
		if err != nil {
			return "", err
		}

		SetSecret(s.UserName, res)
		s.AccessToken = res.AccessToken
	}

	return s.AccessToken, nil
}

// GetRefreshToken ...
func GetRefreshToken() (string, error) {
	// Get Secret Info
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")
	buf, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read secret file: %v", err)
	}

	var s secret
	json.Unmarshal(buf, &s)

	if time.Now().After(s.RefreshTokenExpiresTime) {
		return "", fmt.Errorf("Refresh token was already expired")
	}

	return s.RefreshToken, nil
}

// RemoveSecretFile ...
func RemoveSecretFile() error {
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")
	return os.Remove(secretFile)
}
