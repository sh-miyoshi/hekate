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
func InitConfig(useHTTPS bool, secretKey string) {
	tokenSecretKey = secretKey
	if useHTTPS {
		protoSchema = "https"
	} else {
		protoSchema = "http"
	}
}

func signToken(projectName string, claims jwt.Claims) (string, error) {
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get project")
	}
	switch project.TokenConfig.SigningAlgorithm {
	case "RS256":
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(tokenSecretKey))
	default:
		return "", errors.New("Unexpected Token Signing Algorithm")
	}
}

// GenerateAccessToken ...
func GenerateAccessToken(audiences []string, request Request) (string, error) {
	if tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	// TODO(use Audience in jwt.StandardClaims after merging PR(https://github.com/dgrijalva/jwt-go/pull/355))

	now := time.Now()
	claims := AccessTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    request.Issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			NotBefore: 0,
			Subject:   request.UserID,
		},
		[]string{},
		audiences,
	}

	// Set user roles
	user, err := db.GetInst().UserGet(request.UserID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user")
	}
	for _, role := range user.Roles {
		claims.Roles = append(claims.Roles, role)
	}

	return signToken(request.ProjectName, claims)
}

// GenerateRefreshToken ...
func GenerateRefreshToken(sessionID string, audiences []string, request Request) (string, error) {
	if tokenSecretKey == "" {
		return "", errors.New("Did not initialize config yet")
	}

	// TODO(use Audience in jwt.StandardClaims after merging PR(https://github.com/dgrijalva/jwt-go/pull/355))

	now := time.Now()
	claims := &RefreshTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    request.Issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			NotBefore: 0,
			Subject:   request.UserID,
		},
		sessionID,
		audiences,
	}

	return signToken(request.ProjectName, claims)
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
