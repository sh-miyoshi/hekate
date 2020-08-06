package config

import (
	"flag"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"gopkg.in/yaml.v2"
)

// CheckLoginResDirStruct check a struct of login resource directory
func (c *GlobalConfig) CheckLoginResDirStruct() *errors.Error {
	// TODO(implement this)
	return nil
}

// InitConfig ...
func InitConfig(osArgs []string) (*GlobalConfig, *errors.Error) {
	res := &GlobalConfig{}

	cfile, err := getConfigFileName(osArgs)
	if err != nil {
		return nil, errors.Append(err, "Failed to get parse config file")
	}

	// Set by config file
	if cfile != "" {
		fp, err := os.Open(cfile)
		if err != nil {
			return nil, errors.New("", "Failed to open config file: %v", err)
		}
		defer fp.Close()

		if err := yaml.NewDecoder(fp).Decode(res); err != nil {
			return nil, errors.New("", "Failed to decode config yaml: %v", err)
		}
	}

	// Set by os.Env
	var port, env string
	setEnvVar("HEKATE_ADMIN_NAME", &res.AdminName)
	setEnvVar("HEKATE_ADMIN_PASSWORD", &res.AdminPassword)
	setEnvVar("HEKATE_SERVER_PORT", &port)
	if port != "" {
		res.Port, _ = strconv.Atoi(port)
	}
	setEnvVar("HEKATE_SERVER_BIND_ADDR", &res.BindAddr)
	setEnvVar("HEKATE_SERVER_ENV", &env)
	if strings.ToLower(env) == "debug" {
		res.ModeDebug = true
	}
	setEnvVar("HEKATE_DB_TYPE", &res.DB.Type)
	setEnvVar("HEKATE_DB_CONNECT_STRING", &res.DB.ConnectionString)
	setEnvVar("HEKATE_LOGIN_PAGE_RES", &res.UserLoginResourceDir)

	// Set by command line args

	// "config" flag here is just to avoid an error.
	var c string
	flag.StringVar(&c, "config", "", "config file path")

	flag.StringVar(&res.AdminName, "admin", res.AdminName, "name of administrator")
	flag.StringVar(&res.AdminPassword, "password", res.AdminPassword, "password of administrator")
	flag.IntVar(&res.Port, "port", res.Port, "port number of server")
	flag.StringVar(&res.BindAddr, "bind-addr", res.BindAddr, "bind address of server")
	flag.BoolVar(&res.HTTPSConfig.Enabled, "https", res.HTTPSConfig.Enabled, "start server with https")
	flag.StringVar(&res.HTTPSConfig.CertFile, "https-cert-file", res.HTTPSConfig.CertFile, "cert file path of https")
	flag.StringVar(&res.HTTPSConfig.KeyFile, "https-key-file", res.HTTPSConfig.KeyFile, "key file path of https")
	flag.StringVar(&res.LogFile, "logfile", res.LogFile, "file path for log, output to STDOUT if empty")
	flag.BoolVar(&res.ModeDebug, "debug", res.ModeDebug, "output debug log")
	flag.StringVar(&res.DB.Type, "db-type", res.DB.Type, "type of database")
	flag.StringVar(&res.DB.ConnectionString, "db-conn-str", res.DB.ConnectionString, "database connection string")
	flag.StringVar(&res.UserLoginResourceDir, "login-res", res.UserLoginResourceDir, "directory path for user login")
	flag.Parse()

	// Validate config
	if err := res.Validate(); err != nil {
		return nil, errors.Append(err, "Failed to validate config")
	}

	return res, nil
}

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		*target = val
	}
}

// getConfigFileName return config file name if -config is in os.Args
func getConfigFileName(args []string) (string, *errors.Error) {
	configFilePath := ""
	for i, arg := range args {
		re := regexp.MustCompile(`^--?config=?`)
		if re.MatchString(arg) {
			// arg is one of the following
			//    -config <yaml>
			//    -config=<yaml>
			//   --config <yaml>
			//   --config=<yaml>

			v := strings.Split(arg, "=")
			if len(v) == 1 {
				// arg maybe `-config <yaml>` or `--config <yaml>`
				if i >= len(args)-1 {
					return "", errors.New("", "no config file name")
				}

				nextArg := args[i+1]
				if nextArg[0] == '-' {
					// nextArg is not a config file, but a flag such as "--logfile"
					return "", errors.New("", "nextArg is not a config file name, but a flag")
				}

				configFilePath = nextArg
			} else {
				for j := 1; j < len(v); j++ {
					configFilePath += v[j]
					configFilePath += "=" // split by =
				}
				configFilePath = strings.TrimSuffix(configFilePath, "=") // remove last =
			}
			break
		}
	}
	return configFilePath, nil
}
