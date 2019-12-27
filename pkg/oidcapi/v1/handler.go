package oidc

import (
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"net/http"
)

// ConfigGetHandler method return a configuration of OpenID Connect
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	res := Config{}
	// TODO(set value)

	jwthttp.ResponseWrite(w, "ConfigGetHandler", &res)
}
