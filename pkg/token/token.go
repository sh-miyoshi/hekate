package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"strings"
	"time"
)

var (
	tokenIssuer    string
	tokenSecretKey string
)

type accessTokenClaims struct {
	jwt.StandardClaims

	roles []string
}

type refreshTokenClaims struct {
	jwt.StandardClaims
}

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
	claims := &accessTokenClaims{
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
	user, err := db.GetInst().User.Get(request.ProjectID, request.UserID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user")
	}
	for _, role := range user.Roles {
		claims.roles = append(claims.roles, role)
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
	claims := &refreshTokenClaims{
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

// ValidateToken ...
func ValidateToken(tokenString string) error {
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
