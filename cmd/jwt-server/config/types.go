package config

// DBInfo ...
type DBInfo struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connection_string"`
}

// GlobalConfig ...
type GlobalConfig struct {
	Port           int    `yaml:"server_port"`
	BindAddr       string `yaml:"server_bind_address"`
	LogFile        string `yaml:"logfile"`
	ModeDebug      bool   `yaml:"debug_mode"`
	DB             DBInfo `yaml:"db"`
	AdminName      string `yaml:"admin_name"`
	AdminPassword  string `yaml:"admin_password"`
	TokenIssuer    string `yaml:"token_issuer"`
	TokenSecretKey string `yaml:"token_secret_key"`
}
