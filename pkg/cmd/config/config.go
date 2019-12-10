package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// SystemConfig ...
type SystemConfig struct {
	ServerAddr  string `yaml:"server_addr"`
	EnableDebug bool   `yaml:"enable_debug_log"`
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

	return nil
}

// Get ...
func Get() *SystemConfig {
	return &sysConf
}
