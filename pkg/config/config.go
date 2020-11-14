package config

import (
	"flag"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"gopkg.in/yaml.v2"
)

var inst = GlobalConfig{}

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

	if c.LoginSessionExpiresIn == 0 {
		return errors.New("Invalid config", "login session expires time is 0")
	}

	if c.SSOExpiresIn == 0 {
		return errors.New("Invalid config", "sso expires time is 0")
	}

	if c.DBGCInterval == 0 {
		return errors.New("Invalid config", "interval of db gc is 0")
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

func (c *GlobalConfig) setLoginResource() *errors.Error {
	// directory struct
	// .
	// ├── consent.html  : consent page
	// ├── error.html    : error page
	// ├── index.html    : login page
	// └── static        : directory of static assets

	dir := c.UserLoginResourceDir
	pubMsg := "invalid login resource directory struct"
	c.LoginResource.ConsentPage = path.Join(dir, "/consent.html")
	if _, err := os.Stat(c.LoginResource.ConsentPage); err != nil {
		return errors.New(pubMsg, "Failed to get consent page: %v", err)
	}
	c.LoginResource.ErrorPage = path.Join(dir, "/error.html")
	if _, err := os.Stat(c.LoginResource.ErrorPage); err != nil {
		return errors.New(pubMsg, "Failed to get error page: %v", err)
	}
	c.LoginResource.IndexPage = path.Join(dir, "/index.html")
	if _, err := os.Stat(c.LoginResource.IndexPage); err != nil {
		return errors.New(pubMsg, "Failed to get login page: %v", err)
	}
	c.LoginResource.DeviceLoginPage = path.Join(dir, "/devicelogin.html")
	if _, err := os.Stat(c.LoginResource.DeviceLoginPage); err != nil {
		return errors.New(pubMsg, "Failed to get device login page: %v", err)
	}
	// static directory is option, so does not require check

	return nil
}

// InitConfig ...
func InitConfig(osArgs []string) *errors.Error {
	cfile, err := getConfigFileName(osArgs)
	if err != nil {
		return errors.Append(err, "Failed to get parse config file")
	}

	// Set by config file
	if cfile != "" {
		fp, err := os.Open(cfile)
		if err != nil {
			return errors.New("Broken config", "Failed to open config file: %v", err)
		}
		defer fp.Close()

		if err := yaml.NewDecoder(fp).Decode(&inst); err != nil {
			return errors.New("Broken config", "Failed to decode config yaml: %v", err)
		}
	}

	// Set by os.Env
	var env string
	setEnvVar("HEKATE_ADMIN_NAME", &inst.AdminName)
	setEnvVar("HEKATE_ADMIN_PASSWORD", &inst.AdminPassword)
	if err := setEnvInt("HEKATE_SERVER_PORT", &inst.Port); err != nil {
		return errors.New("Invalid os env", "Failed to get port number: %v", err)
	}
	setEnvVar("HEKATE_SERVER_BIND_ADDR", &inst.BindAddr)
	setEnvVar("HEKATE_SERVER_ENV", &env)
	if strings.ToLower(env) == "debug" {
		inst.ModeDebug = true
	}
	setEnvVar("HEKATE_DB_TYPE", &inst.DB.Type)
	setEnvVar("HEKATE_DB_CONNECT_STRING", &inst.DB.ConnectionString)
	setEnvVar("HEKATE_AUDIT_DB_TYPE", &inst.AuditDB.Type)
	setEnvVar("HEKATE_AUDIT_DB_CONNECT_STRING", &inst.AuditDB.ConnectionString)
	setEnvVar("HEKATE_LOGIN_PAGE_RES", &inst.UserLoginResourceDir)
	if err := setEnvUint("HEKATE_LOGIN_SESSION_EXPIRES_IN", &inst.LoginSessionExpiresIn); err != nil {
		return errors.New("Invalid os env", "Failed to get login session expires time: %v", err)
	}
	if err := setEnvUint("HEKATE_SSO_EXPIRES_IN", &inst.SSOExpiresIn); err != nil {
		return errors.New("Invalid os env", "Failed to get sso expires time: %v", err)
	}
	setEnvVar("HEKATE_LOGIN_PAGE_RES", &inst.UserLoginResourceDir)
	if err := setEnvUint("HEKATE_DBGC_INTERVAL", &inst.DBGCInterval); err != nil {
		return errors.New("Invalid os env", "Failed to get db gc interval: %v", err)
	}

	// Set by command line args

	// "config" flag here is just to avoid an error.
	var c string
	flag.StringVar(&c, "config", "", "config file path")

	flag.StringVar(&inst.AdminName, "admin", inst.AdminName, "name of administrator")
	flag.StringVar(&inst.AdminPassword, "password", inst.AdminPassword, "password of administrator")
	flag.IntVar(&inst.Port, "port", inst.Port, "port number of server")
	flag.StringVar(&inst.BindAddr, "bind-addr", inst.BindAddr, "bind address of server")
	flag.BoolVar(&inst.HTTPSConfig.Enabled, "https", inst.HTTPSConfig.Enabled, "start server with https")
	flag.StringVar(&inst.HTTPSConfig.CertFile, "https-cert-file", inst.HTTPSConfig.CertFile, "cert file path of https")
	flag.StringVar(&inst.HTTPSConfig.KeyFile, "https-key-file", inst.HTTPSConfig.KeyFile, "key file path of https")
	flag.StringVar(&inst.LogFile, "logfile", inst.LogFile, "file path for log, output to STDOUT if empty")
	flag.BoolVar(&inst.ModeDebug, "debug", inst.ModeDebug, "output debug log")
	flag.StringVar(&inst.DB.Type, "db-type", inst.DB.Type, "type of database")
	flag.StringVar(&inst.DB.ConnectionString, "db-conn-str", inst.DB.ConnectionString, "database connection string")
	flag.StringVar(&inst.AuditDB.Type, "audit-db-type", inst.AuditDB.Type, "type of audit events database")
	flag.StringVar(&inst.AuditDB.ConnectionString, "audit-db-conn-str", inst.AuditDB.ConnectionString, "audit database connection string")
	flag.Uint64Var(&inst.LoginSessionExpiresIn, "login-session-expires", inst.LoginSessionExpiresIn, "expires time of login session [sec]")
	flag.Uint64Var(&inst.SSOExpiresIn, "sso-expires", inst.SSOExpiresIn, "expires time of single sign on [sec]")
	flag.StringVar(&inst.UserLoginResourceDir, "login-res", inst.UserLoginResourceDir, "directory path for user login")
	flag.Uint64Var(&inst.DBGCInterval, "dbgc-interval", inst.DBGCInterval, "interval time of garbage collector for expired sessions [sec]")
	flag.Parse()

	// Set supported type
	inst.SupportedResponseType = []string{
		"code",
		"id_token",
		"token",
		"code id_token",
		"code token",
		"id_token token",
		"code id_token token",
		// TODO(support type "none")
	}
	inst.SupportedScore = []string{"openid"}
	inst.LoginStaticResourceURL = "/resource/login"

	// Validate config
	if err := inst.Validate(); err != nil {
		return errors.Append(err, "Failed to validate config")
	}

	if err := inst.setLoginResource(); err != nil {
		return errors.Append(err, "Failed to set login resource")
	}

	return nil
}

// Get ...
func Get() *GlobalConfig {
	return &inst
}

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		*target = val
	}
}

func setEnvInt(key string, target *int) error {
	var tmp string
	setEnvVar(key, &tmp)
	if tmp != "" {
		var err error
		*target, err = strconv.Atoi(tmp)
		return err
	}

	return nil
}

func setEnvUint(key string, target *uint64) error {
	var tmp string
	setEnvVar(key, &tmp)
	if tmp != "" {
		var err error
		*target, err = strconv.ParseUint(tmp, 10, 64)
		return err
	}

	return nil
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
					return "", errors.New("Invalid args", "no config file name")
				}

				nextArg := args[i+1]
				if nextArg[0] == '-' {
					// nextArg is not a config file, but a flag such as "--logfile"
					return "", errors.New("Invalid args", "nextArg is not a config file name, but a flag")
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
