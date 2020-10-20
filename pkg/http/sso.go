package http

import (
	"crypto/x509"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
)

// SetSSOSessionToCookie ...
func SetSSOSessionToCookie(w http.ResponseWriter, projectName, userID, issuer string) *errors.Error {
	prj, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		return errors.Append(err, "Failed to get project config")
	}

	// TODO(set new value to config?)
	lifeSpan := prj.TokenConfig.AccessTokenLifeSpan

	req := token.Request{
		Issuer:      issuer,
		ExpiredTime: time.Second * time.Duration(lifeSpan),
		ProjectName: projectName,
		UserID:      userID,
	}
	tkn, err := token.GenerateSSOToken(req)
	if err != nil {
		return errors.Append(err, "Failed to generate SSO token")
	}

	cookie := &http.Cookie{
		Name:     "HEKATE_LOGIN_SESSION",
		Value:    tkn,
		MaxAge:   int(lifeSpan),
		Secure:   oidc.IsCookieSecure(),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return nil
}

// GetLoginUserIDFromSSOSessionCookie ...
func GetLoginUserIDFromSSOSessionCookie(cookie *http.Cookie, projectName string) (string, *errors.Error) {
	var claims jwt.StandardClaims
	tkn, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		project, err := db.GetInst().ProjectGet(projectName)
		if err != nil {
			return nil, errors.Append(err, "Failed to get project")
		}

		switch token.Method {
		case jwt.SigningMethodRS256:
			key, err := x509.ParsePKCS1PublicKey(project.TokenConfig.SignPublicKey)
			if err != nil {
				return nil, errors.New("Invalid request", "Failed to parse public key: %v", err)
			}
			return key, nil
		}

		return nil, errors.New("Invalid request", "unknown token sigining method")
	})

	if err != nil || !tkn.Valid {
		return "", errors.New("Invalid request", "Token in cookie is not valid: %v", err)
	}

	return claims.Subject, nil
}
