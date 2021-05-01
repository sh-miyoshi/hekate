package mongo

import (
	"time"
)

type tokenConfig struct {
	AccessTokenLifeSpan  uint   `bson:"access_token_life_span"`
	RefreshTokenLifeSpan uint   `bson:"refresh_token_life_span"`
	SigningAlgorithm     string `bson:"signing_algorithm"`
	SignPublicKey        []byte `bson:"sign_public_key"`
	SignSecretKey        []byte `bson:"sign_secret_key"`
}

type passwordPolicy struct {
	MinimumLength       uint     `bson:"length"`
	NotUserName         bool     `bson:"not_user_name"`
	BlackList           []string `bson:"black_list"`
	UseCharacter        string   `bson:"use_character"`
	UseDigit            bool     `bson:"use_digit"`
	UseSpecialCharacter bool     `bson:"use_special_character"`
}

type userLock struct {
	Enabled          bool `bson:"enabled"`
	MaxLoginFailure  uint `bson:"max_login_failure"`
	LockDuration     uint `bson:"lock_duration"`
	FailureResetTime uint `bson:"failure_reset_time"`
}

type projectInfo struct {
	Name            string         `bson:"name"`
	CreatedAt       time.Time      `bson:"create_at"`
	TokenConfig     *tokenConfig   `bson:"token_config"`
	PermitDelete    bool           `bson:"permit_delete"`
	AllowGrantTypes []string       `bson:"allow_grant_types"`
	PasswordPolicy  passwordPolicy `bson:"password_policy"`
	UserLock        userLock       `bson:"user_lock"`
}

type session struct {
	UserID       string    `bson:"user_id"`
	ProjectName  string    `bson:"project_name"`
	SessionID    string    `bson:"session_id"`
	CreatedAt    time.Time `bson:"created_at"`
	ExpiresIn    int64     `bson:"expires_in"`
	FromIP       string    `bson:"from_ip"`
	LastAuthTime time.Time `bson:"last_auth_time"`
}

type loginSession struct {
	SessionID           string    `bson:"session_id"`
	Code                string    `bson:"code"`
	ExpiresDate         time.Time `bson:"expires_in"`
	Scopes              []string  `bson:"scopes"`
	ResponseType        []string  `bson:"response_type"`
	ClientID            string    `bson:"client_id"`
	RedirectURI         string    `bson:"redirect_uri"`
	Nonce               string    `bson:"nonce"`
	ProjectName         string    `bson:"project_name"`
	ResponseMode        string    `bson:"response_mode"`
	Prompt              []string  `bson:"prompt"`
	UserID              string    `bson:"user_id"`
	LoginDate           time.Time `bson:"login_date"`
	CodeChallenge       string    `bson:"code_challenge"`
	CodeChallengeMethod string    `bson:"code_challenge_method"`
}

type lockState struct {
	Locked            bool        `bson:"locked"`
	VerifyFailedTimes []time.Time `bson:"verify_failed_times"`
}

type otpInfo struct {
	ID         string `bson:"id"`
	PrivateKey string `bson:"private_key"`
	Enabled    bool   `bson:"enabled"`
}

type userInfo struct {
	ID           string    `bson:"id"`
	ProjectName  string    `bson:"project_name"`
	Name         string    `bson:"name"`
	EMail        string    `bson:"email"`
	CreatedAt    time.Time `bson:"created_at"`
	PasswordHash string    `bson:"password_hash"`
	SystemRoles  []string  `bson:"system_roles"`
	CustomRoles  []string  `bson:"custom_roles"`
	LockState    lockState `bson:"lock_state"`
	OTPInfo      otpInfo   `bson:"otp_info"`
}

type clientInfo struct {
	ID                  string    `bson:"id"`
	ProjectName         string    `bson:"project_name"`
	Secret              string    `bson:"secret"`
	AccessType          string    `bson:"access_type"`
	CreatedAt           time.Time `bson:"created_at"`
	AllowedCallbackURLs []string  `bson:"allowed_callback_urls"`
}

type customRole struct {
	ID          string    `bson:"id"`
	Name        string    `bson:"name"`
	CreatedAt   time.Time `bson:"created_at"`
	ProjectName string    `bson:"project_name"`
}

type customRoleInUser struct {
	ProjectName  string `bson:"project_name"`
	UserID       string `bson:"user_id"`
	CustomRoleID string `bson:"custom_role_id"`
}

type device struct {
	DeviceCode     string    `bson:"device_code"`
	UserCode       string    `bson:"user_code"`
	ProjectName    string    `bson:"project_name"`
	ExpiresIn      int64     `bson:"expires_in"`
	CreatedAt      time.Time `bson:"created_at"`
	LoginSessionID string    `bson:"login_session_id"`
}
