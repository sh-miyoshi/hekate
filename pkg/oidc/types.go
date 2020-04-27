package oidc

import (
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// AuthRequest ...
type AuthRequest struct {
	// Required
	Scope        string `validate:"required"`
	ResponseType string `validate:"required"`
	ClientID     string `validate:"required"`
	RedirectURI  string `validate:"required,url"`

	// Recommend
	State string

	// Optional
	Nonce  string
	Prompt string
	MaxAge int

	// TODO(implement this)
	// ResponseMode string // response_mode(OPTIONAL)
	// Display string // display(OPTIONAL)
	// UILocales string // ui_locales(OPTIONAL)
	// IDTokenHint string // id_token_hint(OPTIONAL)
	// LoginHint string // login_hint(OPTIONAL)
	// ACRValues string // acr_values(OPTIONAL)
}

// UserLoginInfo ...
type UserLoginInfo struct {
	Scope        string
	ResponseType string
	ClientID     string
	RedirectURI  string
	Nonce        string
	MaxAge       int
}

func validatePrompt(prompts string) error {
	v := strings.Split(prompts, " ")
	if strings.Contains(prompts, "none") && len(v) != 1 {
		return ErrInvalidRequest
	}

	// TODO change response
	for _, prompt := range v {
		switch prompt {
		case "login":
			// login is supported
		case "consent":
			return ErrConsentRequired
		case "select_account":
			return ErrAccountSelectionRequired
		default:
			return ErrInvalidRequest
		}
	}

	return nil
}

// Validate ...
func (r *AuthRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return ErrInvalidRequest
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
		return ErrUnsupportedResponseType
	}

	// Check prompt
	if r.Prompt != "" {
		if err := validatePrompt(r.Prompt); err != nil {
			return err
		}
	}

	if r.MaxAge < 1 {
		return ErrInvalidRequest
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
