package oidc

import (
	"net/url"
)

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
