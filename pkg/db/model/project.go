package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/stretchr/stew/slice"
)

// GrantType ...
type GrantType struct {
	value string
}

// String method returns a name of grant type
func (t GrantType) String() string {
	return t.value
}

// CharacterType ...
type CharacterType string

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  uint
	RefreshTokenLifeSpan uint
	SigningAlgorithm     string
	SignPublicKey        []byte
	SignSecretKey        []byte
}

// PasswordPolicy ...
type PasswordPolicy struct {
	MinimumLength       uint
	NotUserName         bool
	BlackList           []string
	UseCharacter        CharacterType
	UseDigit            bool
	UseSpecialCharacter bool
}

// UserLock ...
type UserLock struct {
	Enabled          bool
	MaxLoginFailure  uint
	LockDuration     uint
	FailureResetTime uint
}

// ProjectInfo ...
type ProjectInfo struct {
	Name            string
	CreatedAt       time.Time
	TokenConfig     *TokenConfig
	PermitDelete    bool
	AllowGrantTypes []GrantType
	PasswordPolicy  PasswordPolicy
	UserLock        UserLock
}

// ProjectFilter ...
type ProjectFilter struct {
	Name string
}

const (
	// DefaultAccessTokenExpiresInSec is default expires time for access token(5 minutes)
	DefaultAccessTokenExpiresInSec = 5 * 60

	// DefaultRefreshTokenExpiresInSec is default expires time for refresh token(14 days)
	DefaultRefreshTokenExpiresInSec = 14 * 24 * 60 * 60

	// DefaultMaxLoginFailure ...
	DefaultMaxLoginFailure = 5

	// DefaultLockDuration is default lock duration(10 minutes)
	DefaultLockDuration = 10 * 60

	// DefaultFailureResetTime is default reset time of login failure(10 minutes)
	DefaultFailureResetTime = 10 * 60
)

var (
	// Defines of Project Error

	// ErrProjectAlreadyExists ...
	ErrProjectAlreadyExists = errors.New("Project already exists", "Project already exists")
	// ErrNoSuchProject ...
	ErrNoSuchProject = errors.New("No such project", "No such project")
	// ErrDeleteBlockedProject ...
	ErrDeleteBlockedProject = errors.New("Project is blocked by delete", "Project cannot be deleted")
	// ErrProjectValidateFailed ...
	ErrProjectValidateFailed = errors.New("Project validation failed", "Project validation failed")

	// Grant Types

	// GrantTypeClientCredentials ...
	GrantTypeClientCredentials = GrantType{"client_credentials"}
	// GrantTypeAuthorizationCode ...
	GrantTypeAuthorizationCode = GrantType{"authorization_code"}
	// GrantTypeRefreshToken ...
	GrantTypeRefreshToken = GrantType{"refresh_token"}
	// GrantTypePassword ...
	GrantTypePassword = GrantType{"password"}

	// Character Types

	// CharacterTypeLower ...
	CharacterTypeLower = CharacterType("lower")
	// CharacterTypeUpper ...
	CharacterTypeUpper = CharacterType("upper")
	// CharacterTypeBoth ...
	CharacterTypeBoth = CharacterType("both")
	// CharacterTypeEither ...
	CharacterTypeEither = CharacterType("either")
	// AllCharacterTypes ...
	AllCharacterTypes = []CharacterType{CharacterTypeLower, CharacterTypeUpper, CharacterTypeBoth, CharacterTypeEither}
)

// ProjectInfoHandler ...
type ProjectInfoHandler interface {
	Add(ent *ProjectInfo) *errors.Error
	Delete(name string) *errors.Error
	GetList(filter *ProjectFilter) ([]*ProjectInfo, *errors.Error)
	Update(ent *ProjectInfo) *errors.Error
}

func (p *PasswordPolicy) validate() *errors.Error {
	if p.UseCharacter != "" && !slice.Contains(AllCharacterTypes, p.UseCharacter) {
		return errors.Append(ErrProjectValidateFailed, "Invalid Character type")
	}
	return nil
}

// Validate ...
func (p *ProjectInfo) Validate() *errors.Error {
	if !ValidateProjectName(p.Name) {
		return errors.Append(ErrProjectValidateFailed, "Invalid Project Name format")
	}

	if !ValidateTokenSigningAlgorithm(p.TokenConfig.SigningAlgorithm) {
		return errors.Append(ErrProjectValidateFailed, "Invalid Token Signing Algorithm")
	}

	if !ValidateLifeSpan(p.TokenConfig.AccessTokenLifeSpan) {
		return errors.Append(ErrProjectValidateFailed, "Access Token Life Span must >= 1")
	}

	if !ValidateLifeSpan(p.TokenConfig.RefreshTokenLifeSpan) {
		return errors.Append(ErrProjectValidateFailed, "Refresh Token Life Span must >= 1")
	}

	if err := p.PasswordPolicy.validate(); err != nil {
		return err
	}

	return nil
}

// GetGrantType ...
func GetGrantType(str string) (GrantType, *errors.Error) {
	if str == GrantTypeClientCredentials.String() {
		return GrantTypeClientCredentials, nil
	}
	if str == GrantTypeAuthorizationCode.String() {
		return GrantTypeAuthorizationCode, nil
	}
	if str == GrantTypeRefreshToken.String() {
		return GrantTypeRefreshToken, nil
	}
	if str == GrantTypePassword.String() {
		return GrantTypePassword, nil
	}

	return GrantType{}, errors.New("No such grant type", "No such grant type")
}
