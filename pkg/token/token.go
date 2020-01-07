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
	protoSchema    string
)

// InitConfig ...
func InitConfig(useHTTPS bool) {
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
		// TODO(fix bug)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(project.TokenConfig.SignSecretKey)
		// token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		// key, err := jwt.ParseRSAPrivateKeyFromPEM(project.TokenConfig.SignSecretKey)
		// if err != nil {
		// 	return "", err
		// }
		// return token.SignedString(key)
	default:
		return "", errors.New("Unexpected Token Signing Algorithm")
	}
}

// GenerateAccessToken ...
func GenerateAccessToken(audiences []string, request Request) (string, error) {
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
		request.ProjectName,
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
		request.ProjectName,
		sessionID,
		audiences,
	}

	return signToken(request.ProjectName, claims)
}

// ValidateAccessToken ...
func ValidateAccessToken(claims *AccessTokenClaims, tokenString string, expectIssuer string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(claims.Project)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get project")
		}

		return project.TokenConfig.SignSecretKey, nil
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
func ValidateRefreshToken(claims *RefreshTokenClaims, tokenString string, expectIssuer string) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(claims.Project)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get project")
		}

		return project.TokenConfig.SignSecretKey, nil
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
