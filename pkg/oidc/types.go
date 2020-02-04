package oidc

import (
	validator "github.com/go-playground/validator/v10"
	"strings"
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

	// Check Response Type
	ok := false
FOR_LABEL:
	for _, typ := range strings.Split(r.ResponseType, " ") {
		for _, support := range GetSupportedResponseType() {
			if typ == support {
				ok = true
				break FOR_LABEL
			}
		}
	}
	if !ok {
		return errors.New("Unsupported response type specified")
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

// Error ...
type Error struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
	Code        int    `json:"status_code"`
}

// Error ...
func (e *Error) Error() string {
	return e.Name
}

// *) error definition is in errors.go
