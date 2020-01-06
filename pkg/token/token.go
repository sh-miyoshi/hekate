package token

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	tokenSecretKey string
	protoSchema    string
)

// InitConfig ...
func InitConfig(secretKey string) {
	tokenSecretKey = secretKey
	protoSchema = "http" // TODO
}

// GenerateAccessToken ...
func GenerateAccessToken(request Request, audiences []string) (string, error) {
	if tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	// TODO(use []string after merging PR(https://github.com/dgrijalva/jwt-go/pull/355))
	aud := "["
	for _, a := range audiences {
		aud += a + ","
	}
	aud = strings.TrimSuffix(aud, ",")
	aud += "]"

	now := time.Now()
	claims := AccessTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    request.Issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			Audience:  aud,
			NotBefore: 0,
			Subject:   request.UserID,
		},
		[]string{},
	}

	// Set user roles
	user, err := db.GetInst().UserGet(request.UserID)
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
func GenerateRefreshToken(sessionID string, request Request) (string, error) {
	if tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	now := time.Now()
	claims := &RefreshTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    request.Issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			Audience:  request.UserID,
			NotBefore: 0,
			Subject:   request.UserID,
		},
		sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecretKey))
}

// ValidateAccessToken ...
func ValidateAccessToken(claims *AccessTokenClaims, tokenString string, expectIssuer string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return errors.Wrap(err, "Failed to parse token")
	}

	if !token.Valid {
		return errors.New("Invalid token is specified")
	}

	// Token Validate
	ti := claims.Issuer
	if len(claims.Issuer) > len(expectIssuer) {
		ti = claims.Issuer[:len(expectIssuer)]
	}
	if ti != expectIssuer {
		logger.Debug("Unexpected token issuer: want %s, got %s", expectIssuer, ti)
		return errors.New("Unexpected token issuer")
	}

	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		return errors.New("Token is expired")
	}

	return nil
}

// ValidateRefreshToken ...
func ValidateRefreshToken(claims *RefreshTokenClaims, tokenString string, issuer string) error {
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

	if claims.Issuer != issuer {
		return errors.New("Unexpected token issuer")
	}

	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		return errors.New("Token is expired")
	}

	return nil
}

// GetFullIssuer ...
func GetFullIssuer(r *http.Request) string {
	re := regexp.MustCompile(`/api/v1/project/[^/]+`)
	url := re.FindString(r.URL.Path)
	res := fmt.Sprintf("%s://%s%s", protoSchema, r.Host, url)
	return strings.TrimSuffix(res, "/")
}

// GetExpectIssuer ...
func GetExpectIssuer(r *http.Request) string {
	return fmt.Sprintf("%s://%s", protoSchema, r.Host)
}
