package oidc

import (
	"sort"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/stretchr/stew/slice"
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
	Nonce               string
	Prompt              []string
	MaxAge              int64
	ResponseMode        string
	IDTokenHint         string
	CodeChallenge       string
	CodeChallengeMethod string

	Request string

	// TODO(implement this)
	// Display string // display(OPTIONAL)
	// UILocales string // ui_locales(OPTIONAL)
	// ACRValues string // acr_values(OPTIONAL)
}

func validatePrompt(prompts []string) *errors.Error {
	for _, prompt := range prompts {
		switch prompt {
		case "login", "select_account", "consent":
			// correct values
		case "none":
			if len(prompts) != 1 {
				return errors.ErrInvalidRequest
			}
		default:
			return errors.ErrInvalidRequest
		}
	}

	return nil
}

func validateResponseType(types, supportedTypes []string) *errors.Error {
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

	if ok := slice.Contains(supportedTypes, s); !ok {
		return errors.ErrUnsupportedResponseType
	}

	return nil
}

func validateResponseMode(mode string) *errors.Error {
	// TODO(add support form_post)
	modes := []string{"query", "fragment"}
	if !slice.Contains(modes, mode) {
		return errors.ErrInvalidRequest
	}
	return nil
}

func validateScope(scope string, supportedScope []string) *errors.Error {
	scopes := strings.Split(scope, " ")
	for _, s := range scopes {
		if !slice.Contains(scopes, s) {
			return errors.ErrInvalidScope
		}
	}

	return nil
}

func validateCodeChallenge(challenge string, method string) *errors.Error {
	if challenge == "" {
		if method != "" {
			return errors.ErrInvalidRequest
		}
	} else {
		methods := []string{"plain", "S256"}
		if !slice.Contains(methods, method) {
			return errors.ErrInvalidRequest
		}
	}
	return nil
}

// Validate ...
func (r *AuthRequest) Validate() *errors.Error {
	if err := validator.New().Struct(r); err != nil {
		return errors.Append(errors.ErrInvalidRequest, err.Error())
	}

	// TODO(support request)
	if r.Request != "" {
		return errors.ErrRequestNotSupported
	}

	cfg := config.Get()

	// Check Scope
	if err := validateScope(r.Scope, cfg.SupportedScope); err != nil {
		return errors.Append(err, "Failed to validate scope %v", r.Scope)
	}

	// Check Response Type
	if err := validateResponseType(r.ResponseType, cfg.SupportedResponseType); err != nil {
		return errors.Append(err, "Failed to validate response type %v", r.ResponseType)
	}

	// Check prompt
	if err := validatePrompt(r.Prompt); err != nil {
		return errors.Append(err, "Failed to validate prompt %v", r.Prompt)
	}

	// Check Response mode
	if err := validateResponseMode(r.ResponseMode); err != nil {
		return errors.Append(err, "Failed to validate response mode %s", r.ResponseMode)
	}

	// Check CodeChallengeMethod
	if err := validateCodeChallenge(r.CodeChallenge, r.CodeChallengeMethod); err != nil {
		return errors.Append(err, "Failed to validate code challenge %s with method %s", r.CodeChallenge, r.CodeChallengeMethod)
	}

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
