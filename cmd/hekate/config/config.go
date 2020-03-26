package config

import (
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		*target = val
	}
}

// InitConfig ...
func InitConfig(filePath string) (*GlobalConfig, error) {
	res := &GlobalConfig{}

	fp, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open config file")
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(res); err != nil {
		return nil, errors.Wrap(err, "Failed to decode config yaml")
	}

	setEnvVar("HEKATE_ADMIN_NAME", &res.AdminName)
	setEnvVar("HEKATE_ADMIN_PASSWORD", &res.AdminPassword)

	return res, nil
}
