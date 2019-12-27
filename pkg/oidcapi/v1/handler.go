package oidc

import (
	"fmt"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"net/http"
	"strings"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// ConfigGetHandler method return a configuration of OpenID Connect
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(set correct ProtoSchema default: https)
	addr := strings.TrimSuffix(fmt.Sprintf("http://%s%s", r.Host, r.RequestURI), "/.well-known/openid-configuration")
	logger.Debug("Issuer: %s", addr)

	res := Config{
		Issuer:                           addr,
		AuthorizationEndpoint:            addr + "/openid-connect/auth",
		TokenEndpoint:                    addr + "/openid-connect/token",
		UserinfoEndpoint:                 addr + "/openid-connect/userinfo",
		JwksURI:                          addr + "/openid-connect/certs",
		ScopesSupported:                  []string{"openid"},
		ResponseTypesSupported:           []string{}, // TODO(set value)
		SubjectTypesSupported:            []string{}, // TODO(set value)
		IDTokenSigningAlgValuesSupported: []string{"HS256"},
		ClaimsSupported: []string{
			"iss",
			"aud",
			"sub",
			"exp",
			"jti",
			"iat",
			"nbf",
		},
	}

	jwthttp.ResponseWrite(w, "ConfigGetHandler", &res)
}

// TokenHandler ...
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// internal server error?
		return
	}

	logger.Info("Form: %v", r.Form)

	res := TokenResponse{}
	// TODO(set value)

	jwthttp.ResponseWrite(w, "TokenHandler", &res)
}