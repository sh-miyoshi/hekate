package oidc

import (
	"crypto/x509"
	"encoding/json"
	"github.com/dvsekhvalnov/jose2go/base64url"
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
	issuer := token.GetFullIssuer(r)
	logger.Debug("Issuer: %s", issuer)

	res := Config{
		Issuer:                 issuer,
		AuthorizationEndpoint:  issuer + "/openid-connect/auth",
		TokenEndpoint:          issuer + "/openid-connect/token",
		UserinfoEndpoint:       issuer + "/openid-connect/userinfo",
		JwksURI:                issuer + "/openid-connect/certs",
		ScopesSupported:        []string{"openid"},
		ResponseTypesSupported: []string{},         // TODO(set value)
		SubjectTypesSupported:  []string{"public"}, // TODO(set value)
		IDTokenSigningAlgValuesSupported: []string{
			"RS256",
		},
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

		audiences := []string{user.ID}
		clientID := r.Form.Get("client_id")
		if clientID != "" {
			audiences = append(audiences, clientID)
		}

		accessTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
			ProjectName: user.ProjectName,
			UserID:      user.ID,
		}

		res.AccessToken, err = token.GenerateAccessToken(audiences, accessTokenReq)
		if err != nil {
			logger.Error("Failed to get JWT token: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		refreshTokenReq := token.Request{
			Issuer:      token.GetFullIssuer(r),
			ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
			ProjectName: user.ProjectName,
			UserID:      user.ID,
		}

		sessionID := uuid.New().String()
		res.RefreshToken, err = token.GenerateRefreshToken(sessionID, audiences, refreshTokenReq)
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
	case "refresh_token":
		// TODO(implement token get by refresh_token)
		logger.Error("Not Implemented yet")
		writeTokenErrorResponse(w)
		return
	}

	logger.Info("No such Grant Type: %s", r.Form.Get("grant_type"))
	writeTokenErrorResponse(w)
}

// CertsHandler ...
func CertsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	jwk := JWKInfo{
		KeyID:        uuid.New().String(),
		Algorithm:    project.TokenConfig.SigningAlgorithm,
		PublicKeyUse: "sig",
	}

	switch jwk.Algorithm {
	case "RS256":
		jwk.KeyType = "RSA"
		key, err := x509.ParsePKCS1PublicKey(project.TokenConfig.SignPublicKey)
		if err != nil {
			logger.Error("Failed to parse RSA public key: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		e := util.Int2bytes(uint64(key.E))
		jwk.E = base64url.Encode(e)
		jwk.N = base64url.Encode(key.N.Bytes())
	}

	res := JWKSet{}
	res.Keys = append(res.Keys, jwk)
	jwthttp.ResponseWrite(w, "CertsHandler", &res)
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
