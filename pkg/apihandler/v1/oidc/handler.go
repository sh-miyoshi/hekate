package oidc

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/oidc"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/user"
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
		SubjectTypesSupported:  []string{"public"},
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
	if errors.Cause(err) == model.ErrNoSuchProject {
		http.Error(w, "Project Not Found", http.StatusNotFound)
		return
	}

	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	if err := oidc.ClientAuth(clientID, clientSecret); err != nil {
		if errors.Cause(err) == oidc.ErrClientAuthFailed {
			logger.Info("Failed to authenticate client: %s", clientID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			logger.Error("Failed to authenticate client: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	var tkn *oidc.TokenResponse

	// Authetication
	switch r.Form.Get("grant_type") {
	case "password":
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")
		tkn, err = oidc.ReqAuthByPassword(project, uname, passwd, r)
	case "refresh_token":
		refreshToken := r.Form.Get("refresh_token")
		tkn, err = oidc.ReqAuthByRefreshToken(project, refreshToken, r)
	case "authorization_code":
		codeID := r.Form.Get("code")
		tkn, err = oidc.ReqAuthByCode(project, codeID, r)
	default:
		logger.Info("No such Grant Type: %s", r.Form.Get("grant_type"))
		writeTokenErrorResponse(w)
		return
	}

	if err != nil {
		if errors.Cause(err) == oidc.ErrRequestVerifyFailed {
			logger.Info("Failed to verify request: %v", err)
			writeTokenErrorResponse(w)
		} else {
			logger.Error("Failed to verify request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

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
}

// CertsHandler ...
func CertsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	res, err := oidc.GenerateJWKSet(project.TokenConfig.SigningAlgorithm, project.TokenConfig.SignPublicKey)
	if err != nil {
		logger.Error("Failed to generate JWT set: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwthttp.ResponseWrite(w, "CertsHandler", res)
}

// AuthGETHandler ...
func AuthGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	authReq := oidc.NewAuthRequest(queries)
	if err := authReq.Validate(); err != nil {
		logger.Info("Failed to validate request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(switch by request)

	// return end user auth prompt
	code := oidc.RegisterUserLoginSession(authReq)
	oidc.WriteUserLoginPage(code, projectName, w)
}

// AuthPOSTHandler ...
func AuthPOSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

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

	// TODO(switch by request)

	// return end user auth prompt
	code := oidc.RegisterUserLoginSession(authReq)
	oidc.WriteUserLoginPage(code, projectName, w)
}

// UserLoginHandler ...
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Debug("Form: %v", r.Form)

	// Verify user login session code
	authReq, err := oidc.UserLoginVerify(r.Form.Get("login_verify_code"))
	if err != nil {
		logger.Info("Failed to verify user login session: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := user.Verify(projectName, uname, passwd)
	if err != nil {
		if errors.Cause(err) == user.ErrAuthFailed {
			logger.Info("Failed to authenticate user %s: %v", uname, err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to verify user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	code, _ := oidc.GenerateAuthCode(authReq.ClientID, authReq.RedirectURI, usr.ID)
	values := url.Values{}
	values.Set("code", code)
	if authReq.State != "" {
		values.Set("state", authReq.State)
	}

	req, err := http.NewRequest("GET", authReq.RedirectURI, nil)
	if err != nil {
		logger.Error("Failed to create response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

// UserInfoHandler ...
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := jwthttp.ValidateAPIRequest(r)
	if err != nil {
		logger.Info("Failed to validate header: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := db.GetInst().UserGet(claims.Subject)
	if err != nil {
		// If token validate accepted, user absolutely exists
		logger.Error("Failed to get user: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := &UserInfo{
		Subject:  claims.Subject,
		UserName: user.Name,
	}

	jwthttp.ResponseWrite(w, "UserInfoHandler", res)
}
