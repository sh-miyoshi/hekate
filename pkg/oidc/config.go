package oidc

var (
	expiresTimeSec   uint64
	userLoginHTML    string
	userLoginResPath string
)

// InitConfig ...
func InitConfig(authCodeExpiresTimeSec uint64, authCodeUserLoginFile string, authCodeUserLoginResPath string) {
	expiresTimeSec = authCodeExpiresTimeSec
	userLoginHTML = authCodeUserLoginFile
	userLoginResPath = authCodeUserLoginResPath
}

// GetSupportedResponseType ...
func GetSupportedResponseType() []string {
	return []string{
		"code",
	}
}
