package output

import (
	"encoding/json"
	"fmt"

	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
)

// ConfigFormat ...
type ConfigFormat struct {
	config *config.SystemConfig
}

// NewConfigFormat ...
func NewConfigFormat(config *config.SystemConfig) *ConfigFormat {
	return &ConfigFormat{
		config: config,
	}
}

// ToText ...
func (f *ConfigFormat) ToText() (string, error) {
	res := fmt.Sprintf("Server Address:  %s\n", f.config.ServerAddr)
	res += fmt.Sprintf("Default Project: %s\n", f.config.DefaultProject)
	res += fmt.Sprintf("Client ID:       %s\n", f.config.ClientID)
	res += fmt.Sprintf("Client Secret:   %s\n", f.config.ClientSecret)
	res += fmt.Sprintf("TLS Skip Verify: %v\n", f.config.Insecure)
	res += fmt.Sprintf("Request Timeout: %d", f.config.RequestTimeout)
	return res, nil
}

// ToJSON ...
func (f *ConfigFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.config)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
