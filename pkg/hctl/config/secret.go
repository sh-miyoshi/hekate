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
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
)

// Secret ...
type Secret struct {
	ProjectName             string    `json:"projectName"`
	UserName                string    `json:"userName"`
	AccessToken             string    `json:"accessToken"`
	AccessTokenExpiresDate  time.Time `json:"accessTokenExpiresDate"`
	RefreshToken            string    `json:"refreshToken"`
	RefreshTokenExpiresDate time.Time `json:"refreshTokenExpiresDate"`
}

// SetSecret ...
func SetSecret(projectName string, userName string, token *oidcapi.TokenResponse) {
	secretFile := filepath.Join(configDir, "secret")

	// doNext
	v := Secret{
		ProjectName:             projectName,
		UserName:                userName,
		AccessToken:             token.AccessToken,
		AccessTokenExpiresDate:  time.Now().Add(time.Second * time.Duration(token.ExpiresIn)),
		RefreshToken:            token.RefreshToken,
		RefreshTokenExpiresDate: time.Now().Add(time.Second * time.Duration(token.RefreshExpiresIn)),
	}

	bytes, _ := json.MarshalIndent(v, "", "  ")
	ioutil.WriteFile(secretFile, bytes, os.ModePerm)
}

// GetSecret ...
func GetSecret() (*Secret, error) {
	// Get Secret Info
	secretFile := filepath.Join(configDir, "secret")
	buf, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read secret file: %v\nYou need to `hctl login` at first", err)
	}

	var s Secret
	json.Unmarshal(buf, &s)
	return &s, nil
}

// GetAccessToken ...
func GetAccessToken() (string, error) {
	s, err := GetSecret()
	if err != nil {
		return "", err
	}

	if time.Now().After(s.RefreshTokenExpiresDate) {
		return "", fmt.Errorf("Token is expired\nPlease run `hctl login`")
	}

	if time.Now().After(s.AccessTokenExpiresDate) {
		// Refresh token by using refresh-token
		res, err := login.DoWithRefresh(s.RefreshToken, login.Info{
			ServerAddr:   sysConf.ServerAddr,
			ProjectName:  s.ProjectName,
			ClientID:     sysConf.ClientID,
			ClientSecret: sysConf.ClientSecret,
			Insecure:     sysConf.Insecure,
			Timeout:      sysConf.RequestTimeout,
		})
		if err != nil {
			return "", err
		}

		SetSecret(s.ProjectName, s.UserName, res)
		s.AccessToken = res.AccessToken
	}

	return s.AccessToken, nil
}

// RemoveSecretFile ...
func RemoveSecretFile() error {
	secretFile := filepath.Join(configDir, "secret")
	print.Debug("Removing secret file: %s", secretFile)
	return os.Remove(secretFile)
}
