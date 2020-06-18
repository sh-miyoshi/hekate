package config

import (
	"io/ioutil"
	"os"
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

func setDefaultParams() {
	sysConf.ServerAddr = "http://localhost:18443"
	sysConf.DefaultProject = "master"
	sysConf.ClientID = "portal"
}

// InitConfig ...
func InitConfig(confDir string) error {
	if confDir == "" {
		// set default path
		conf, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		confDir = filepath.Join(conf, "hekate")
	}

	// Read Config File
	fname := filepath.Join(confDir, "config.yaml")
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		// if error is no such file, create file with default params
		if os.IsNotExist(err) {
			setDefaultParams()
			if err := SaveToFile(); err != nil {
				return err
			}
		} else {
			return err
		}
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

// SaveToFile ...
func SaveToFile() error {
	// TODO(implement this)
	return nil
}
