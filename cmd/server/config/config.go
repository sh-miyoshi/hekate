package config

import (
	yaml "gopkg.in/yaml.v2"
	"os"
)

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		target = &val
	}
}

// InitConfig ...
func InitConfig(filePath string) (*GlobalConfig, error) {
	res := &GlobalConfig{}

	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(res); err != nil {
		return nil, err
	}

	setEnvVar("JWT_SERVER_ADMIN_NAME", &res.AdminName)
	setEnvVar("JWT_SERVER_ADMIN_PASSWORD", &res.AdminPassword)

	return res, nil
}
