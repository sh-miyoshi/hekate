package http

import (
	"net/http"
	"strings"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

func parseHTTPHeaderToken(tokenString string) (string, *errors.Error) {
	var splitToken []string
	if strings.Contains(tokenString, "bearer ") {
		splitToken = strings.Split(tokenString, "bearer ")
	} else if strings.Contains(tokenString, "Bearer ") {
		splitToken = strings.Split(tokenString, "Bearer ")
	} else {
		return "", errors.New("Invalid request", "token format is missing")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

// ValidateAPIRequest ...
func ValidateAPIRequest(req *http.Request) (*token.AccessTokenClaims, *errors.Error) {
	auth, ok := req.Header["Authorization"]
	if !ok || len(auth) != 1 {
		return nil, errors.New("Failed to get Authorization header", "Failed to get Authorization header")
	}
	tokenString, err := parseHTTPHeaderToken(auth[0])
	if err != nil {
		return nil, errors.Append(err, "Failed to get token from header")
	}
	claims := &token.AccessTokenClaims{}
	issuer := token.GetExpectIssuer(req)
	if err := token.ValidateAccessToken(claims, tokenString, issuer); err != nil {
		return nil, errors.Append(err, "Failed to validate token")
	}
	return claims, nil
}

// Authorize ...
func Authorize(req *http.Request, projectName string, reqTrgRes role.Resource, reqRoleType role.Type) *errors.Error {
	claims, err := ValidateAPIRequest(req)
	if err != nil {
		return errors.Append(err, "Failed to validate token")
	}

	// return ok when request has cluster-role
	if role.Authorize(claims.ResourceAccess.SystemManagement.Roles, role.ResCluster, reqRoleType) {
		return nil
	}

	// Authorize API Request
	if !role.Authorize(claims.ResourceAccess.SystemManagement.Roles, reqTrgRes, reqRoleType) {
		return errors.New("Do not have permission", "Do not have permission")
	}

	// check project
	if claims.Project != projectName {
		return errors.New("Wrong project", "Wrong project")
	}

	return nil
}
