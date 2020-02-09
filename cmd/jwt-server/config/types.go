package config

// DBInfo ...
type DBInfo struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

// HTTPSConfig ...
type HTTPSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert-file"`
	KeyFile  string `yaml:"key-file"`
}

// GlobalConfig ...
type GlobalConfig struct {
	Port                 int         `yaml:"server_port"`
	BindAddr             string      `yaml:"server_bind_address"`
	LogFile              string      `yaml:"logfile"`
	ModeDebug            bool        `yaml:"debug_mode"`
	DB                   DBInfo      `yaml:"db"`
	AdminName            string      `yaml:"admin_name"`
	AdminPassword        string      `yaml:"admin_password"`
	AuthCodeExpiresTime  uint64      `yaml:"oidc_auth_code_expires_time"`
	UserLoginPage        string      `yaml:"oidc_user_login_page"`
	UserLoginResourceDir string      `yaml:"oidc_user_login_resource_dir"`
	HTTPSConfig          HTTPSConfig `yaml:"https"`
}
