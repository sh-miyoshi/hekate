package oidc

import (
	"net/url"
	"strconv"
	"strings"
)

// NewAuthRequest ...
func NewAuthRequest(values url.Values) *AuthRequest {
	request := values.Get("request")
	// TODO(parse request and set params)

	maxAge, _ := strconv.Atoi(values.Get("max_age"))
	prompt := []string{}
	if values.Get("prompt") != "" {
		prompt = strings.Split(values.Get("prompt"), " ")
	}
	responseTypes := strings.Split(values.Get("response_type"), " ")

	resMode := values.Get("response_mode")
	if resMode == "" {
		resMode = "fragment"
		if len(responseTypes) == 1 {
			if responseTypes[0] == "code" || responseTypes[0] == "none" {
				resMode = "query"
			}
		}
	}

	return &AuthRequest{
		Scope:        values.Get("scope"),
		ResponseType: responseTypes,
		ClientID:     values.Get("client_id"),
		RedirectURI:  values.Get("redirect_uri"),
		State:        values.Get("state"),
		Nonce:        values.Get("nonce"),
		Prompt:       prompt,
		MaxAge:       uint(maxAge),
		LoginHint:    values.Get("login_hint"),
		ResponseMode: resMode,
		Request:      request,
	}
}
