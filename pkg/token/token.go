package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"strings"
	"time"
)

// TODO(use secure key)
const testSecretKey = "hjsoeirjhgrokwrejrhoir4owjpjkrghpljrwoi34j"
const issuer = "jwt-server"

func validate(claims jwt.Claims, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(testSecretKey), nil
	})

	if err != nil {
		return err
	}

	logger.Debug("Claims: %v\n", claims)

	if token.Valid {
		return nil
	}
	return fmt.Errorf("Failed to validate token")
}

// ParseHTTPHeaderToken return jwt token from http header
func ParseHTTPHeaderToken(tokenString string) (string, error) {
	var splitToken []string
	if strings.Contains(tokenString, "bearer") {
		splitToken = strings.Split(tokenString, "bearer")
	} else if strings.Contains(tokenString, "Bearer") {
		splitToken = strings.Split(tokenString, "Bearer")
	} else {
		return "", fmt.Errorf("token format is missing")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

// Generate returns jwt token for user
func Generate() (string, error) {
	// TODO(set args)

	claims := &jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // Expired at 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(testSecretKey))
}

// Authenticate validates token
func Authenticate(reqToken string) error {
	claims := jwt.StandardClaims{}
	err := validate(&claims, reqToken)
	if err != nil {
		logger.Info("Failed to auth token %v", err)
		return err
	}
	logger.Debug("claims in token: %v", claims)

	// Validate claims
	if claims.Issuer != issuer {
		logger.Info("Issuer want %s, but got %s", issuer, claims.Issuer)
		return fmt.Errorf("Issuer want %s, but got %s", issuer, claims.Issuer)
	}

	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		logger.Info("Token is expired at %d. now: %d", claims.ExpiresAt, now)
		return fmt.Errorf("Token is expired at %d. now: %d", claims.ExpiresAt, now)
	}

	return nil
}
