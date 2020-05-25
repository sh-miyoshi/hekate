package oidc

import (
	"net/url"
	"strconv"
	"strings"
)

// NewAuthRequest ...
func NewAuthRequest(values url.Values) *AuthRequest {
	maxAge, _ := strconv.Atoi(values.Get("max_age"))
	return &AuthRequest{
		Scope:        values.Get("scope"),
		ResponseType: strings.Split(values.Get("response_type"), " "),
		ClientID:     values.Get("client_id"),
		RedirectURI:  values.Get("redirect_uri"),
		State:        values.Get("state"),
		Nonce:        values.Get("nonce"),
		Prompt:       values.Get("prompt"),
		MaxAge:       uint(maxAge),
	}
}
