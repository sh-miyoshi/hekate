package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

// Request ...
type Request struct {
	ExpiredTime time.Duration
	ProjectID   string
	UserID      string
}

// AccessTokenClaims ...
type AccessTokenClaims struct {
	jwt.StandardClaims

	Roles []string `json:"roles"`
}

// RefreshTokenClaims ...
type RefreshTokenClaims struct {
	jwt.StandardClaims
}
