package oidc

import (
	validator "github.com/go-playground/validator/v10"
	"net/url"
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

// NewAuthRequest ...
func NewAuthRequest(values url.Values) *AuthRequest {
	return &AuthRequest{
		Scope:        values.Get("scope"),
		ResponseType: values.Get("response_type"),
		ClientID:     values.Get("client_id"),
		RedirectURI:  values.Get("redirect_uri"),
		State:        values.Get("state"),
	}
}

// Validate ...
func (r *AuthRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	// TODO(add more validation)

	return nil
}
