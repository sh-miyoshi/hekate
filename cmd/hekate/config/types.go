package config

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// DBInfo ...
type DBInfo struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

// HTTPSConfig ...
type HTTPSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert-file"`
	KeyFile  string `yaml:"key-file"`
}

// GlobalConfig ...
type GlobalConfig struct {
	Port                 int         `yaml:"server_port"`
	BindAddr             string      `yaml:"server_bind_address"`
	LogFile              string      `yaml:"logfile"`
	ModeDebug            bool        `yaml:"debug_mode"`
	DB                   DBInfo      `yaml:"db"`
	AdminName            string      `yaml:"admin_name"`
	AdminPassword        string      `yaml:"admin_password"`
	AuthCodeExpiresTime  uint64      `yaml:"oidc_auth_code_expires_time"`
	UserLoginResourceDir string      `yaml:"oidc_user_login_page_res"`
	HTTPSConfig          HTTPSConfig `yaml:"https"`
	AuditDB              DBInfo      `yaml:"audit_db"`
}

// Validate ...
func (c *GlobalConfig) Validate() *errors.Error {
	if c.Port == 0 || c.Port > 65535 {
		return errors.New("Invalid config", "port number %d is not valid", c.Port)
	}

	if c.AdminName == "" {
		return errors.New("Invalid config", "admin name is empty")
	}

	if c.AdminPassword == "" {
		return errors.New("Invalid config", "admin password is empty")
	}

	if c.AuthCodeExpiresTime == 0 {
		return errors.New("Invalid config", "login session expires time is 0")
	}

	finfo, err := os.Stat(c.UserLoginResourceDir)
	if err != nil {
		return errors.New("Invalid config", "Failed to get login resource info: %v", err)
	}
	if !finfo.IsDir() {
		return errors.New("Invalid config", "login resource path %s is not directory", c.UserLoginResourceDir)
	}

	return nil
}
