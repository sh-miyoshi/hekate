package config

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		*target = val
	}
}

// InitConfig ...
func InitConfig(filePath string) (*GlobalConfig, *errors.Error) {
	res := &GlobalConfig{}

	fp, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("", "Failed to open config file: %v", err)
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(res); err != nil {
		return nil, errors.New("", "Failed to decode config yaml: %v", err)
	}

	setEnvVar("HEKATE_ADMIN_NAME", &res.AdminName)
	setEnvVar("HEKATE_ADMIN_PASSWORD", &res.AdminPassword)
	setEnvVar("HEKATE_DB_TYPE", &res.DB.Type)
	setEnvVar("HEKATE_DB_CONNECT_STRING", &res.DB.ConnectionString)

	return res, nil
}
