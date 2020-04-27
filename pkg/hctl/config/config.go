package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// SystemConfig ...
type SystemConfig struct {
	ConfigDir      string
	EnableDebug    bool
	ServerAddr     string `yaml:"server_addr"`
	DefaultProject string `yaml:"default_project"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
}

var sysConf SystemConfig

// InitConfig ...
func InitConfig(confDir string) error {
	// Read Config File
	fname := filepath.Join(confDir, "config.yaml")
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &sysConf); err != nil {
		return err
	}

	sysConf.ConfigDir = confDir
	sysConf.EnableDebug = false

	return nil
}

// EnableDebugMode ...
func EnableDebugMode() {
	sysConf.EnableDebug = true
}

// Get ...
func Get() *SystemConfig {
	return &sysConf
}
