package oidc

var (
	expiresTimeSec uint64
	userLoginHTML  string
)

// InitConfig ...
func InitConfig(authCodeExpiresTimeSec uint64, authCodeUserLoginFile string) {
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginHTML = authCodeUserLoginFile
}

// GetSupportedResponseType ...
func GetSupportedResponseType() []string {
	return []string{
		"code",
	}
}
