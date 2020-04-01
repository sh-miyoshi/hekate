package model

import (
	"time"

	"github.com/pkg/errors"
)

// GrantType ...
type GrantType struct {
	value string
}

// String method returns a name of grant type
func (t GrantType) String() string {
	return t.value
}

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  uint
	RefreshTokenLifeSpan uint
	SigningAlgorithm     string
	SignPublicKey        []byte
	SignSecretKey        []byte
}

// ProjectInfo ...
type ProjectInfo struct {
	Name            string
	CreatedAt       time.Time
	TokenConfig     *TokenConfig
	PermitDelete    bool
	AllowGrantTypes []GrantType
}

const (
	// DefaultAccessTokenExpiresTimeSec is default expires time for access token(5 minutes)
	DefaultAccessTokenExpiresTimeSec = 5 * 60

	// DefaultRefreshTokenExpiresTimeSec is default expires time for refresh token(14 days)
	DefaultRefreshTokenExpiresTimeSec = 14 * 24 * 60 * 60
)

var (
	// Defines of Project Error

	// ErrProjectAlreadyExists ...
	ErrProjectAlreadyExists = errors.New("Project Already Exists")
	// ErrNoSuchProject ...
	ErrNoSuchProject = errors.New("No such project")
	// ErrDeleteBlockedProject ...
	ErrDeleteBlockedProject = errors.New("Projects cannot be deleted")
	// ErrProjectValidateFailed ...
	ErrProjectValidateFailed = errors.New("Project Validation Failed")

	// Grant Types

	// GrantTypeClientCredentials ...
	GrantTypeClientCredentials = GrantType{"Client Credentials"}
	// GrantTypeAuthorizationCode ...
	GrantTypeAuthorizationCode = GrantType{"Authorization Code"}
	// GrantTypeRefreshToken ...
	GrantTypeRefreshToken = GrantType{"Refresh Token"}
	// GrantTypePassword ...
	GrantTypePassword = GrantType{"Password"}
)

// ProjectInfoHandler ...
type ProjectInfoHandler interface {
	Add(ent *ProjectInfo) error
	Delete(name string) error
	GetList() ([]*ProjectInfo, error)
	Get(name string) (*ProjectInfo, error)

	// Update method updates existing project
	// It must return error if project is not found
	Update(ent *ProjectInfo) error
}

// Validate ...
func (p *ProjectInfo) Validate() error {
	if !ValidateProjectName(p.Name) {
		return errors.Wrap(ErrProjectValidateFailed, "Invalid Project Name format")
	}

	if !ValidateTokenSigningAlgorithm(p.TokenConfig.SigningAlgorithm) {
		return errors.Wrap(ErrProjectValidateFailed, "Invalid Token Signing Algorithm")
	}

	if !ValidateLifeSpan(p.TokenConfig.AccessTokenLifeSpan) {
		return errors.Wrap(ErrProjectValidateFailed, "Access Token Life Span must >= 1")
	}

	if !ValidateLifeSpan(p.TokenConfig.RefreshTokenLifeSpan) {
		return errors.Wrap(ErrProjectValidateFailed, "Refresh Token Life Span must >= 1")
	}

	return nil
}

// GetGrantType ...
func GetGrantType(str string) (GrantType, error) {
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

	return GrantType{}, errors.New("No such grant type")
}
