package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Request ...
type Request struct {
	Issuer      string
	ExpiredTime time.Duration
	ProjectName string
	UserID      string
	Nonce       string
}

// RoleValue ...
type RoleValue struct {
	Roles []string `json:"roles"`
}

// RoleSet ...
type RoleSet struct {
	SystemManagement RoleValue `json:"system_management"`
	User             RoleValue `json:"user"`
}

// AccessTokenClaims ...
type AccessTokenClaims struct {
	jwt.StandardClaims

	Project        string   `json:"project"`
	Audience       []string `json:"aud"`
	ResourceAccess RoleSet  `json:"resource_access"`
	UserName       string   `json:"preferred_username"`
}

// RefreshTokenClaims ...
type RefreshTokenClaims struct {
	jwt.StandardClaims

	Project   string   `json:"project"`
	SessionID string   `json:"sessionID"`
	Audience  []string `json:"aud"`
}

// IDTokenClaims ...
type IDTokenClaims struct {
	jwt.StandardClaims

	Audience []string `json:"aud"`
	Nonce    string   `json:"nonce"`
	// TODO(auth_time, acr, amr, azp)
	// ref. https://openid-foundation-japan.github.io/openid-connect-core-1_0.ja.html#IDToken
}
