package oidc

import (
	"net/url"
)

// AuthRequest ...
type AuthRequest struct {
	Scope        string // scope(REQUIRED)
	ResponseType string // response_type(REQUIRED)
	ClientID     string // client_id(REQUIRED)
	RedirectURI  string // redirect_uri(REQUIRED)
	State        string // state(RECOMMENDED)

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
	// TODO(add validation)
	// refs.
	//   https://qiita.com/itkr/items/9b4e8d8c6d574137443c
	//   https://github.com/go-playground/validator
	//   https://godoc.org/gopkg.in/go-playground/validator.v9
	return nil
}
