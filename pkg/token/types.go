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
}

// AccessTokenClaims ...
type AccessTokenClaims struct {
	jwt.StandardClaims

	Project  string   `json:"project"`
	Roles    []string `json:"roles"`
	Audience []string `json:"aud"`
}

// RefreshTokenClaims ...
type RefreshTokenClaims struct {
	jwt.StandardClaims

	Project   string   `json:"project"`
	SessionID string   `json:"sessionID"`
	Audience  []string `json:"aud"`
}
