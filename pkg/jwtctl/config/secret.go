package config

import (
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/login"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
	"os"
)

type secret struct {
	UserName string `json:"userName"`
	AccessToken string `json:"accessToken"`
	AccessTokenExpiresTime time.Time `json:"accessTokenExpiresTime"`
	RefreshToken string `json:"refreshToken"`
	RefreshTokenExpiresTime time.Time `json:"refreshTokenExpiresTime"`
}

// SetSecret ...
func SetSecret(userName string,token *tokenapi.TokenResponse) {
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")

	v := secret{
		UserName: userName,
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

	if time.Now().After(s.RefreshTokenExpiresTime) {
		return "", fmt.Errorf("Token is expired\nPlease run `jwtctl login`")
	}

	if time.Now().After(s.AccessTokenExpiresTime) {
		// Refresh token by using refresh-token
		req := tokenapi.TokenRequest{
			Name:     s.UserName,
			Secret:   s.RefreshToken,
			AuthType: "refresh",
		}

		res,err := login.Do(sysConf.ServerAddr, sysConf.ProjectName, &req)
		if err != nil {
			return "", err
		}

		SetSecret(s.UserName, res)
	}

	return s.AccessToken, nil
}