package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"strings"
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

// Validate ...
func Validate(tokenString string) error {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Cause(fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return errors.Wrap(err, "Failed to parse token")
	}

	if !token.Valid {
		return errors.New("Invalid token is specifyed")
	}

	if claims.Issuer != tokenIssuer {
		return errors.New("Unexpected token issuer")
	}

	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		return errors.New("Token is expired")
	}

	return nil
}

// ParseHTTPHeaderToken return jwt token from http header
func ParseHTTPHeaderToken(tokenString string) (string, error) {
	var splitToken []string
	if strings.Contains(tokenString, "bearer") {
		splitToken = strings.Split(tokenString, "bearer")
	} else if strings.Contains(tokenString, "Bearer") {
		splitToken = strings.Split(tokenString, "Bearer")
	} else {
		return "", errors.New("token format is missing")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}
