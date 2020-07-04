package token

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

var (
	protoSchema string
)

// InitConfig ...
func InitConfig(useHTTPS bool) {
	if useHTTPS {
		protoSchema = "https"
	} else {
		protoSchema = "http"
	}
}

func signToken(projectName string, claims jwt.Claims) (string, *errors.Error) {
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		return "", errors.Append(err, "Failed to get project")
	}
	switch project.TokenConfig.SigningAlgorithm {
	case "RS256":
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		key, err := x509.ParsePKCS1PrivateKey(project.TokenConfig.SignSecretKey)
		if err != nil {
			return "", errors.New("Failed to parse private key: %v", err)
		}
		str, err := token.SignedString(key)
		if err != nil {
			return "", errors.New("Failed to signing token: %v", err)
		}
		return str, nil
	default:
		return "", errors.New("Unexpected Token Signing Algorithm")
	}
}

// GenerateAccessToken ...
func GenerateAccessToken(audiences []string, request Request) (string, *errors.Error) {
	user, err := db.GetInst().UserGet(request.ProjectName, request.UserID)
	if err != nil {
		return "", errors.Append(err, "Failed to get user")
	}

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
		audiences,
		RoleSet{
			SystemManagement: RoleValue{
				Roles: []string{},
			},
		},
		user.Name,
	}

	for _, role := range user.SystemRoles {
		claims.ResourceAccess.SystemManagement.Roles = append(claims.ResourceAccess.SystemManagement.Roles, role)
	}
	for _, rid := range user.CustomRoles {
		role, err := db.GetInst().CustomRoleGet(request.ProjectName, rid)
		if err != nil {
			return "", errors.Append(err, "Failed to get custom role name")
		}
		claims.ResourceAccess.User.Roles = append(claims.ResourceAccess.User.Roles, role.Name)
	}

	return signToken(request.ProjectName, claims)
}

// GenerateRefreshToken ...
func GenerateRefreshToken(sessionID string, audiences []string, request Request) (string, *errors.Error) {
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

// GenerateIDToken ...
func GenerateIDToken(audiences []string, request Request) (string, *errors.Error) {
	now := time.Now()
	claims := &IDTokenClaims{
		jwt.StandardClaims{
			Id:        uuid.New().String(),
			Issuer:    request.Issuer,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(request.ExpiredTime).Unix(),
			NotBefore: 0,
			Subject:   request.UserID,
		},
		audiences,
		request.Nonce,
		request.EndUserAuthTime.Unix(),
	}

	return signToken(request.ProjectName, claims)
}

// ValidateAccessToken ...
func ValidateAccessToken(claims *AccessTokenClaims, tokenString string, expectIssuer string) *errors.Error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(claims.Project)
		if err != nil {
			return nil, errors.Append(err, "Failed to get project")
		}

		switch token.Method {
		case jwt.SigningMethodRS256:
			key, err := x509.ParsePKCS1PublicKey(project.TokenConfig.SignPublicKey)
			if err != nil {
				return nil, errors.New("Failed to parse public key: %v", err)
			}
			return key, nil
		}
		return nil, errors.New("unknown token sigining method")
	})

	if err != nil {
		e := err.(*errors.Error)
		return errors.Append(e, "Failed to parse token")
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
func ValidateRefreshToken(claims *RefreshTokenClaims, tokenString string, expectIssuer string) *errors.Error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(claims.Project)
		if err != nil {
			return nil, errors.Append(err, "Failed to get project")
		}

		switch token.Method {
		case jwt.SigningMethodRS256:
			key, err := x509.ParsePKCS1PublicKey(project.TokenConfig.SignPublicKey)
			if err != nil {
				return nil, errors.New("Failed to parse public key: %v", err)
			}
			return key, nil
		}
		return nil, errors.New("unknown token sigining method")
	})

	if err != nil {
		e := err.(*errors.Error)
		return errors.Append(e, "Failed to parse token")
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
