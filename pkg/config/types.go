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

// LoginResource ...
type LoginResource struct {
	IndexPage   string
	ErrorPage   string
	ConsentPage string
}

// GlobalConfig ...
type GlobalConfig struct {
	Port                    int         `yaml:"server_port"`
	BindAddr                string      `yaml:"server_bind_address"`
	LogFile                 string      `yaml:"logfile"`
	ModeDebug               bool        `yaml:"debug_mode"`
	DB                      DBInfo      `yaml:"db"`
	AdminName               string      `yaml:"admin_name"`
	AdminPassword           string      `yaml:"admin_password"`
	LoginSessionExpiresTime uint64      `yaml:"login_session_expires_time"`
	UserLoginResourceDir    string      `yaml:"oidc_user_login_page_res"`
	HTTPSConfig             HTTPSConfig `yaml:"https"`
	AuditDB                 DBInfo      `yaml:"audit_db"`
	DBGCInterval            uint64      `yaml:"dbgc_interval"`
	SSOExpiresTime          uint64      `yaml:"sso_expires_time"`

	SupportedResponseType []string
	SupportedScore        []string
	LoginResource         LoginResource
}
