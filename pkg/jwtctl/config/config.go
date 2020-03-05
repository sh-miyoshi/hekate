package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// SystemConfig ...
type SystemConfig struct {
	ConfigDir   string
	ServerAddr  string `yaml:"server_addr"`
	EnableDebug bool   `yaml:"enable_debug_log"`
	ProjectName string `yaml:"project_name"`
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
