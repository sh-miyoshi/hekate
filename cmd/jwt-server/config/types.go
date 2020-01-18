package config

// DBInfo ...
type DBInfo struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

// GlobalConfig ...
type GlobalConfig struct {
	Port                  int    `yaml:"server_port"`
	BindAddr              string `yaml:"server_bind_address"`
	LogFile               string `yaml:"logfile"`
	ModeDebug             bool   `yaml:"debug_mode"`
	DB                    DBInfo `yaml:"db"`
	AdminName             string `yaml:"admin_name"`
	AdminPassword         string `yaml:"admin_password"`
	AuthCodeExpiresTime   uint64 `yaml:"oidc_auth_code_expires_time"`
	AuthCodeUserLoginFile string `yaml:"oidc_auth_code_user_login_html"`
}
