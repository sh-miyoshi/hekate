package oidc

import (
	"net/url"
	"strconv"
	"strings"
)

// NewAuthRequest ...
func NewAuthRequest(values url.Values) *AuthRequest {
	maxAge, _ := strconv.Atoi(values.Get("max_age"))
	prompt := []string{}
	if values.Get("prompt") != "" {
		prompt = strings.Split(values.Get("prompt"), " ")
	}

	return &AuthRequest{
		Scope:        values.Get("scope"),
		ResponseType: strings.Split(values.Get("response_type"), " "),
		ClientID:     values.Get("client_id"),
		RedirectURI:  values.Get("redirect_uri"),
		State:        values.Get("state"),
		Nonce:        values.Get("nonce"),
		Prompt:       prompt,
		MaxAge:       uint(maxAge),
	}
}
