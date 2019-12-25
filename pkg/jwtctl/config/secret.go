package config

import (
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// GetSecretToken ...
func GetSecretToken() (*tokenapi.TokenResponse, error) {
	// Get Secret Info
	secretFile := filepath.Join(sysConf.ConfigDir, "secret")
	buf, err := ioutil.ReadFile(secretFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read secret file: %v\nYou need to `jwtctl login` at first", err)
	}

	var secret tokenapi.TokenResponse
	json.Unmarshal(buf, &secret)

	// TODO(Validate secret)

	return &secret, nil
}