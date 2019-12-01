package http

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"net/http"
	"strings"
)

func parseHTTPHeaderToken(tokenString string) (string, error) {
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

// ValidateAPIRequest ...
func ValidateAPIRequest(header http.Header) (*token.AccessTokenClaims, error) {
	auth, ok := header["Authorization"]
	if !ok || len(auth) != 1 {
		return nil, errors.New("Failed to get Authorization header")
	}
	tokenString, err := parseHTTPHeaderToken(auth[0])
	if err != nil {
		return nil, errors.New("Failed to get token from header")
	}
	claims := &token.AccessTokenClaims{}
	if err := token.ValidateAccessToken(claims, tokenString); err != nil {
		return nil, errors.Wrap(err, "Failed to validate token")
	}
	return claims, nil
}
