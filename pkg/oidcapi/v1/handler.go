package oidc

import (
	"fmt"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"net/http"
	"strings"
)

// ConfigGetHandler method return a configuration of OpenID Connect
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(set correct ProtoSchema default: https)
	addr := strings.TrimSuffix(fmt.Sprintf("http://%s%s", r.Host, r.RequestURI), "/.well-known/openid-configuration")
	res := Config{
		Issuer:                           addr,
		AuthorizationEndpoint:            addr + "/auth",
		TokenEndpoint:                    addr + "/token",
		UserinfoEndpoint:                 addr + "/userinfo",
		JwksURI:                          addr + "/certs",
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
