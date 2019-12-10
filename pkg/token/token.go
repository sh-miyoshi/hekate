package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
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

// GenerateAccessToken ...
func GenerateAccessToken(request Request) (string, error) {
	if tokenIssuer == "" || tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	now := time.Now()
	claims := AccessTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    tokenIssuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			Audience:  request.UserID,
			NotBefore: 0,
			Subject:   request.UserID,
		},
		[]string{},
	}

	// Set user roles
	user, err := db.GetInst().User.Get(request.ProjectName, request.UserID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user")
	}
	for _, role := range user.Roles {
		claims.Roles = append(claims.Roles, role)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecretKey))
}

// GenerateRefreshToken ...
func GenerateRefreshToken(request Request) (string, error) {
	if tokenIssuer == "" || tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	now := time.Now()
	claims := &RefreshTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    tokenIssuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			Audience:  request.UserID,
			NotBefore: 0,
			Subject:   request.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecretKey))
}

// ValidateAccessToken ...
func ValidateAccessToken(claims *AccessTokenClaims, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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

// ValidateRefreshToken ...
func ValidateRefreshToken(claims *RefreshTokenClaims, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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
