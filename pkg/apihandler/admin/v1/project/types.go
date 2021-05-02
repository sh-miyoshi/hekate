package projectapi

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  uint   `json:"access_token_life_span"`
	RefreshTokenLifeSpan uint   `json:"refresh_token_life_span"`
	SigningAlgorithm     string `json:"signing_algorithm"`
}

// PasswordPolicy ...
type PasswordPolicy struct {
	MinimumLength       uint     `json:"length"`
	NotUserName         bool     `json:"not_user_name"`
	BlackList           []string `json:"black_list"`
	UseCharacter        string   `json:"use_character"`
	UseDigit            bool     `json:"use_digit"`
	UseSpecialCharacter bool     `json:"use_special_character"`
}

// UserLock ...
type UserLock struct {
	Enabled          bool `json:"enabled"`
	MaxLoginFailure  uint `json:"max_logi_failure"`
	LockDuration     uint `json:"lock_duration"`
	FailureResetTime uint `json:"failure_reset_time"`
}

// ProjectCreateRequest ...
type ProjectCreateRequest struct {
	Name            string         `json:"name"`
	TokenConfig     TokenConfig    `json:"token_config"`
	PasswordPolicy  PasswordPolicy `json:"password_policy"`
	AllowGrantTypes []string       `json:"allow_grant_types"`
	UserLock        UserLock       `json:"user_lock"`
}

// ProjectGetResponse ...
type ProjectGetResponse struct {
	Name            string         `json:"name"`
	CreatedAt       string         `json:"created_at"`
	TokenConfig     TokenConfig    `json:"token_config"`
	PasswordPolicy  PasswordPolicy `json:"password_policy"`
	AllowGrantTypes []string       `json:"allow_grant_types"`
	UserLock        UserLock       `json:"user_lock"`
}

// ProjectPutRequest ...
type ProjectPutRequest struct {
	TokenConfig     TokenConfig    `json:"token_config"`
	PasswordPolicy  PasswordPolicy `json:"password_policy"`
	AllowGrantTypes []string       `json:"allow_grant_types"`
	UserLock        UserLock       `json:"user_lock"`
}
