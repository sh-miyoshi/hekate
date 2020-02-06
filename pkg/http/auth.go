package http

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/oidc/token"
	"github.com/sh-miyoshi/jwt-server/pkg/role"
)

func parseHTTPHeaderToken(tokenString string) (string, error) {
	var splitToken []string
	if strings.Contains(tokenString, "bearer ") {
		splitToken = strings.Split(tokenString, "bearer ")
	} else if strings.Contains(tokenString, "Bearer ") {
		splitToken = strings.Split(tokenString, "Bearer ")
	} else {
		return "", errors.New("token format is missing")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

// ValidateAPIRequest ...
func ValidateAPIRequest(req *http.Request) (*token.AccessTokenClaims, error) {
	auth, ok := req.Header["Authorization"]
	if !ok || len(auth) != 1 {
		return nil, errors.New("Failed to get Authorization header")
	}
	tokenString, err := parseHTTPHeaderToken(auth[0])
	if err != nil {
		return nil, errors.New("Failed to get token from header")
	}
	claims := &token.AccessTokenClaims{}
	issuer := token.GetExpectIssuer(req)
	if err := token.ValidateAccessToken(claims, tokenString, issuer); err != nil {
		return nil, errors.Wrap(err, "Failed to validate token")
	}
	return claims, nil
}

// AuthHeader ...
func AuthHeader(req *http.Request, reqTrgRes role.Resource, reqRoleType role.Type) error {
	claims, err := ValidateAPIRequest(req)
	if err != nil {
		return errors.Wrap(err, "Failed to validate token")
	}

	// Authorize API Request
	if !role.GetInst().Authorize(claims.ResourceAccess.SystemManagement.Roles, reqTrgRes, reqRoleType) {
		return errors.New("Do not have authority")
	}

	return nil
}
