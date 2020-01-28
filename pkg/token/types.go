package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

// Request ...
type Request struct {
	Issuer      string
	ExpiredTime time.Duration
	ProjectName string
	UserID      string
	// TODO(add user name)
}

// RoleValue ...
type RoleValue struct {
	Roles []string `json:"roles"`
}

// RoleSet ...
type RoleSet struct {
	SystemManagement RoleValue `json:"system_management"`
}

// AccessTokenClaims ...
type AccessTokenClaims struct {
	jwt.StandardClaims

	Project        string   `json:"project"`
	Audience       []string `json:"aud"`
	ResourceAccess RoleSet  `json:"resource_access"`
}

// RefreshTokenClaims ...
type RefreshTokenClaims struct {
	jwt.StandardClaims

	Project   string   `json:"project"`
	SessionID string   `json:"sessionID"`
	Audience  []string `json:"aud"`
}
