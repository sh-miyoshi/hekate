package oidc

import (
	"sort"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

// AuthRequest ...
type AuthRequest struct {
	// Required
	Scope        string   `validate:"required"`
	ResponseType []string `validate:"required"`
	ClientID     string   `validate:"required"`
	RedirectURI  string   `validate:"required,url"`

	// Recommend
	State string

	// Optional
	Nonce        string
	Prompt       string
	MaxAge       uint
	ResponseMode string

	// TODO(implement this)
	// Display string // display(OPTIONAL)
	// UILocales string // ui_locales(OPTIONAL)
	// IDTokenHint string // id_token_hint(OPTIONAL)
	// LoginHint string // login_hint(OPTIONAL)
	// ACRValues string // acr_values(OPTIONAL)
}

// AuthCodeSession ...
type AuthCodeSession struct {
	Scope        string
	ResponseType []string
	ClientID     string
	RedirectURI  string
	Nonce        string
	MaxAge       uint
	ResponseMode string
	Prompt       string
}

func validatePrompt(prompts string) error {
	v := strings.Split(prompts, " ")
	if strings.Contains(prompts, "none") && len(v) != 1 {
		return ErrInvalidRequest
	}

	for _, prompt := range v {
		switch prompt {
		case "login", "select_account", "consent":
			// TODO(implement this)
		case "none":
			return ErrInteractionRequired
		default:
			return ErrInvalidRequest
		}
	}

	return nil
}

func validateResponseType(types, supportedTypes []string) error {
	// sort types
	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})

	// make string
	s := ""
	for _, typ := range types {
		s += typ + " "
	}
	s = strings.TrimSuffix(s, " ")

	// include check
	for _, support := range supportedTypes {
		if s == support {
			return nil
		}
	}
	return ErrUnsupportedResponseType
}

// Validate ...
func (r *AuthRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return ErrInvalidRequest
	}

	// Check Response Type
	supportedTypes := GetSupportedResponseType()
	if err := validateResponseType(r.ResponseType, supportedTypes); err != nil {
		return err
	}

	// Check prompt
	if r.Prompt != "" {
		if err := validatePrompt(r.Prompt); err != nil {
			return err
		}
	}

	// Check Response mode
	if r.ResponseMode != "" {
		// TODO(add support form_post)
		modes := []string{"query", "fragment"}
		ok := false
		for _, m := range modes {
			if r.ResponseMode == m {
				ok = true
				break
			}
		}

		// TODO return err when query && response_type is not none or code

		if !ok {
			return ErrInvalidRequest
		}
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
