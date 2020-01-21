package oidc

// Config ...
type Config struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported                  []string `json:"claims_supported"`
}

// TokenResponse ...
type TokenResponse struct {
	TokenType        string `json:"token_type"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        uint   `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn uint   `json:"refresh_expires_in"`
	IDToken          string `json:"id_token"`
}

// TokenErrorResponse ...
type TokenErrorResponse struct {
	Error string `json:"error"`
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

// UserInfo ...
type UserInfo struct {
	Subject string `json:"sub"`
	UserName string `json:"preferred_username"`
}