package oidc

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// AuthRequest ...
type AuthRequest struct {
	Scope        string `validate:"required"`
	ResponseType string `validate:"required"`
	ClientID     string `validate:"required"`
	RedirectURI  string `validate:"required,url"`
	State        string

	// TODO(implement this)
	// ResponseMode string // response_mode(OPTIONAL)
	// Nonce string // nonce(OPTIONAL)
	// Display string // display(OPTIONAL)
	// Prompt string // prompt(OPTIONAL)
	// MaxAge string // max_age(OPTIONAL)
	// UILocales string // ui_locales(OPTIONAL)
	// IDTokenHint string // id_token_hint(OPTIONAL)
	// LoginHint string // login_hint(OPTIONAL)
	// ACRValues string // acr_values(OPTIONAL)
}

// Validate ...
func (r *AuthRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	// TODO(add more validation)

	return nil
}

// JWKInfo is a struct for JSON Web Key(JWK) format defined in https://tools.ietf.org/html/rfc7517
type JWKInfo struct {
	KeyType      string `json:"kty"`
	KeyID        string `json:"kid"`
	Algorithm    string `json:"alg"`
	PublicKeyUse string `json:"use"`
	N            string `json:"n,omitempty"` // Use in RSA
	E            string `json:"e,omitempty"` // Use in RSA
	X            string `json:"x,omitempty"` // Use in EC
	Y            string `json:"y,omitempty"` // Use in EC
}

// JWKSet ...
type JWKSet struct {
	Keys []JWKInfo `json:"keys"`
}

// TokenResponse ...
type TokenResponse struct {
	TokenType        string
	AccessToken      string
	ExpiresIn        uint
	RefreshToken     string
	RefreshExpiresIn uint
	IDToken          string
}

var (
	// ErrClientAuthFailed ...
	ErrClientAuthFailed = errors.New("client authentication failed")
)
