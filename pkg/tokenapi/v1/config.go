package tokenapi

type adminConfig struct {
	ID       string
	Password string
}

var config adminConfig

// InitAdminConfig initialize admin info for tokenapi
func InitAdminConfig(id string, password string) {
	config.ID = id
	config.Password = password
}
