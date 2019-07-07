package token

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

type tokenConfig struct {
	ExpiredTime time.Duration
	Issuer      string
	SecretKey   string
}

var config tokenConfig

// InitConfig initialize config of token package
func InitConfig(expiredTime time.Duration, issuer string, secretKey string) {
	config.ExpiredTime = expiredTime
	config.Issuer = issuer
	config.SecretKey = secretKey
}

func validate(claims jwt.Claims, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.SecretKey), nil
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
	claims := &jwt.StandardClaims{
		Issuer:    config.Issuer,
		ExpiresAt: time.Now().Add(config.ExpiredTime).Unix(), // Expired at 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.SecretKey))
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
	if claims.Issuer != config.Issuer {
		logger.Info("Issuer want %s, but got %s", config.Issuer, claims.Issuer)
		return fmt.Errorf("Issuer want %s, but got %s", config.Issuer, claims.Issuer)
	}

	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		logger.Info("Token is expired at %d. now: %d", claims.ExpiresAt, now)
		return fmt.Errorf("Token is expired at %d. now: %d", claims.ExpiresAt, now)
	}

	return nil
}
