package oidc

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"time"
)

// ConfigGetHandler method return a configuration of OpenID Connect
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	issuer := token.GetIssuer(r)
	logger.Debug("Issuer: %s", issuer)

	res := Config{
		Issuer:                           issuer,
		AuthorizationEndpoint:            issuer + "/openid-connect/auth",
		TokenEndpoint:                    issuer + "/openid-connect/token",
		UserinfoEndpoint:                 issuer + "/openid-connect/userinfo",
		JwksURI:                          issuer + "/openid-connect/certs",
		ScopesSupported:                  []string{"openid"},
		ResponseTypesSupported:           []string{}, // TODO(set value)
		SubjectTypesSupported:            []string{}, // TODO(set value)
		IDTokenSigningAlgValuesSupported: []string{"HS256"},
		ClaimsSupported: []string{
			"iss",
			"aud",
			"sub",
			"exp",
			"jti",
			"iat",
			"nbf",
		},
	}

	jwthttp.ResponseWrite(w, "ConfigGetHandler", &res)
}

// TokenHandler ...
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	if err := r.ParseForm(); err != nil {
		// TODO internal server error?
		return
	}

	logger.Info("Form: %v", r.Form)
	switch r.Form.Get("grant_type") {
	case "password":
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")

		user, err := db.GetInst().UserGetByName(projectName, uname)
		if err != nil {
			if err == model.ErrNoSuchUser {
				logger.Info("No such user %s in project %s", user.Name, projectName)
				writeTokenErrorResponse(w)
			} else {
				logger.Error("Failed to get user id: %+v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		hash := util.CreateHash(passwd)
		if user.PasswordHash != hash {
			logger.Info("password authentication failed")
			writeTokenErrorResponse(w)
			return
		}

		// Get Project Info for Token Config
		project, err := db.GetInst().ProjectGet(projectName)
		if err == model.ErrNoSuchProject {
			http.Error(w, "Project Not Found", http.StatusNotFound)
			return
		}

		// Generate JWT Token
		res := TokenResponse{
			TokenType:        "Bearer",
			ExpiresIn:        project.TokenConfig.AccessTokenLifeSpan,
			RefreshExpiresIn: project.TokenConfig.RefreshTokenLifeSpan,
		}

		accessTokenReq := token.Request{
			Issuer:      token.GetIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
			ProjectName: user.ProjectName,
			UserID:      user.ID,
		}

		res.AccessToken, err = token.GenerateAccessToken(accessTokenReq)
		if err != nil {
			logger.Error("Failed to get JWT token: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		refreshTokenReq := token.Request{
			Issuer:      token.GetIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
			ProjectName: user.ProjectName,
			UserID:      user.ID,
		}

		sessionID := uuid.New().String()
		res.RefreshToken, err = token.GenerateRefreshToken(sessionID, refreshTokenReq)
		if err != nil {
			logger.Error("Failed to get JWT token: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := db.GetInst().NewSession(user.ID, sessionID, res.RefreshExpiresIn, r.RemoteAddr); err != nil {
			logger.Error("Failed to register refresh token session token: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		jwthttp.ResponseWrite(w, "TokenHandler", &res)
		return
	}

	logger.Info("No such Grant Type: %s", r.Form.Get("grant_type"))
	writeTokenErrorResponse(w)
}

func writeTokenErrorResponse(w http.ResponseWriter) {
	res := TokenErrorResponse{
		Error: "invalid_request",
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Pragma", "no-cache")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a token error response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
