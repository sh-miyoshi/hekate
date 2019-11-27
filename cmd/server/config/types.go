package config

// GlobalConfig ...
type GlobalConfig struct {
	Port           int    `yaml:"server_port"`
	BindAddr       string `yaml:"server_bind_address"`
	LogFile        string `yaml:"logfile"`
	ModeDebug      bool   `yaml:"debug_mode"`
	AdminName      string `yaml:"admin_name"`
	AdminPassword  string `yaml:"admin_password"`
	TokenIssuer    string `yaml:"token_issuer"`
	TokenSecretKey string `yaml:"token_secret_key"`
}
