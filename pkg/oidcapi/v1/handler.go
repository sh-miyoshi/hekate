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
	"github.com/sh-miyoshi/jwt-server/pkg/oidc"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"net/url"
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
		ResponseTypesSupported: []string{"code"},
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
		logger.Info("Failed to parse form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Debug("Form: %v", r.Form)

	// Get Project Info for Token Config
	project, err := db.GetInst().ProjectGet(projectName)
	if err == model.ErrNoSuchProject {
		http.Error(w, "Project Not Found", http.StatusNotFound)
		return
	}

	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	if err := oidc.ClientAuth(clientID, clientSecret); err != nil {
		// TODO(internal server error)
		logger.Info("Failed to authenticate client: %s", clientID)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tkn *oidc.TokenResponse
	var statusCode int
	var message string

	// Authetication
	switch r.Form.Get("grant_type") {
	case "password":
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")

		tkn, statusCode, message = oidc.AuthByPassword(project, uname, passwd, r)
	case "refresh_token":
		refreshToken := r.Form.Get("refresh_token")
		tkn, statusCode, message = oidc.AuthByRefreshToken(project, refreshToken, r)
	case "authorization_code":
		// Validate code
		codeID := r.Form.Get("code")
		tkn, statusCode, message = oidc.AuthByCode(project, codeID, r)
	default:
		logger.Info("No such Grant Type: %s", r.Form.Get("grant_type"))
		writeTokenErrorResponse(w)
		return
	}

	switch statusCode {
	case http.StatusInternalServerError:
		logger.Error(message)
		http.Error(w, "Internal Server Error", statusCode)
	case http.StatusBadRequest:
		logger.Info(message)
		writeTokenErrorResponse(w)
	case http.StatusOK:
		res := &TokenResponse{
			TokenType:        tkn.TokenType,
			AccessToken:      tkn.AccessToken,
			ExpiresIn:        tkn.ExpiresIn,
			RefreshToken:     tkn.RefreshToken,
			RefreshExpiresIn: tkn.RefreshExpiresIn,
			IDToken:          tkn.IDToken,
		}

		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Pragma", "no-cache")
		jwthttp.ResponseWrite(w, "TokenHandler", res)
	default:
		logger.Error("Program Bug: code %d is not implemented", statusCode)
		http.Error(w, "Internal Server Error", statusCode)
	}
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

// AuthGETHandler ...
func AuthGETHandler(w http.ResponseWriter, r *http.Request) {
	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	authReq := oidc.NewAuthRequest(queries)
	if err := authReq.Validate(); err != nil {
		logger.Info("Failed to validate request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// return end user auth prompt
	// TODO(set header)
	oidc.WriteUserLoginPage(w)

	// // Debug(following code is temporary code)
	// // TODO(set correct user id)
	// users, _ := db.GetInst().UserGetList("master")
	// code, _ := oidc.GenerateAuthCode(authReq.ClientID, authReq.RedirectURI, users[0])
	// values := url.Values{}
	// values.Set("code", code)
	// if authReq.State != "" {
	// 	values.Set("state", authReq.State)
	// }

	// req, err := http.NewRequest("GET", authReq.RedirectURI, nil)
	// if err != nil {
	// 	logger.Error("Failed to create response: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.URL.RawQuery = values.Encode()

	// http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

// AuthPOSTHandler ...
func AuthPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Debug("Form: %v", r.Form)

	authReq := oidc.NewAuthRequest(r.Form)
	if err := authReq.Validate(); err != nil {
		logger.Info("Failed to validate request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		// TODO(return correct error response)
		return
	}

	// TODO(return end user auth prompt)

	// Debug(following code is temporary code)
	// TODO(set correct user id)
	users, _ := db.GetInst().UserGetList("master")
	code, _ := oidc.GenerateAuthCode(authReq.ClientID, authReq.RedirectURI, users[0])
	values := url.Values{}
	values.Set("code", code)
	if authReq.State != "" {
		values.Set("state", authReq.State)
	}

	req, err := http.NewRequest("POST", authReq.RedirectURI, nil)
	if err != nil {
		logger.Error("Failed to create response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

func writeTokenErrorResponse(w http.ResponseWriter) {
	res := TokenErrorResponse{
		Error: "invalid_request",
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a token error response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
