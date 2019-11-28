package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

var (
	tokenIssuer    string
	tokenSecretKey string
)

// InitConfig ...
func InitConfig(issuer string, secretKey string) {
	tokenIssuer = issuer
	tokenSecretKey = secretKey
}

// Generate ...
func Generate(request Request) (string, error) {
	if tokenIssuer == "" || tokenSecretKey == "" {
		return "", errors.Cause(fmt.Errorf("Did not initialize config yet"))
	}

	claims := &jwt.StandardClaims{
		Issuer:    tokenIssuer,
		ExpiresAt: time.Now().Add(request.ExpiredTime).Unix(),
		Audience:  request.Audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecretKey))
}
