package projectapi

import (
	"time"
)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  uint   `json:"accessTokenLifeSpan"`
	RefreshTokenLifeSpan uint   `json:"refreshTokenLifeSpan"`
	SigningAlgorithm     string `json:"signingAlgorithm"`
}

// PasswordPolicy ...
type PasswordPolicy struct {
	MinimumLength       uint     `json:"length"`
	NotUserName         bool     `json:"notUserName"`
	BlackList           []string `json:"blackList"`
	UseCharacter        string   `json:"useCharacter"`
	UseDigit            bool     `json:"useDigit"`
	UseSpecialCharacter bool     `json:"useSpecialCharacter"`
}

// UserLock ...
type UserLock struct {
	Enabled          bool          `json:"enabled"`
	MaxLoginFailure  uint          `json:"maxLoginFailure"`
	LockDuration     time.Duration `json:"lockDuration"`
	FailureResetTime time.Duration `json:"failureResetTime"`
}

// ProjectCreateRequest ...
type ProjectCreateRequest struct {
	Name            string         `json:"name"`
	TokenConfig     TokenConfig    `json:"tokenConfig"`
	PasswordPolicy  PasswordPolicy `json:"passwordPolicy"`
	AllowGrantTypes []string       `json:"allowGrantTypes"`
	UserLock        UserLock       `json:"userLock"`
}

// ProjectGetResponse ...
type ProjectGetResponse struct {
	Name            string         `json:"name"`
	CreatedAt       time.Time      `json:"createdAt"`
	TokenConfig     TokenConfig    `json:"tokenConfig"`
	PasswordPolicy  PasswordPolicy `json:"passwordPolicy"`
	AllowGrantTypes []string       `json:"allowGrantTypes"`
	UserLock        UserLock       `json:"userLock"`
}

// ProjectPutRequest ...
type ProjectPutRequest struct {
	TokenConfig     TokenConfig    `json:"tokenConfig"`
	PasswordPolicy  PasswordPolicy `json:"passwordPolicy"`
	AllowGrantTypes []string       `json:"allowGrantTypes"`
	UserLock        UserLock       `json:"userLock"`
}
