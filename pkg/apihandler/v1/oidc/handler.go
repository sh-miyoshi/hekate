package oidc

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/oidc"
	"github.com/sh-miyoshi/jwt-server/pkg/oidc/token"
	"github.com/sh-miyoshi/jwt-server/pkg/user"
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
		ResponseTypesSupported: oidc.GetSupportedResponseType(),
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
	state := r.Form.Get("state")

	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		writeTokenErrorResponse(w, oidc.ErrInvalidRequestObject, state)
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
		if errors.Cause(err) == oidc.ErrInvalidClient {
			logger.Info("Failed to authenticate client: %s", clientID)
			writeTokenErrorResponse(w, oidc.ErrInvalidClient, state)
		} else {
			logger.Error("Failed to authenticate client: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, state)
		}
		return
	}

	var tkn *oidc.TokenResponse

	// TODO(consider redirect_uri)

	// Authetication
	switch r.Form.Get("grant_type") {
	case "password":
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")
		tkn, err = oidc.ReqAuthByPassword(project, uname, passwd, r)
	case "refresh_token":
		refreshToken := r.Form.Get("refresh_token")
		tkn, err = oidc.ReqAuthByRefreshToken(project, clientID, refreshToken, r)

		if err != nil && errors.Cause(err) == model.ErrNoSuchSession {
			logger.Info("Refresh token is already revoked")
			writeTokenErrorResponse(w, oidc.ErrInvalidRequest, state)
			return
		}
	case "authorization_code":
		codeID := r.Form.Get("code")
		tkn, err = oidc.ReqAuthByCode(project, clientID, codeID, r)
	default:
		logger.Info("No such Grant Type: %s", r.Form.Get("grant_type"))
		writeTokenErrorResponse(w, oidc.ErrInvalidGrant, state)
		return
	}

	if err != nil {
		if errors.Cause(err) == oidc.ErrInvalidRequest {
			logger.Info("Failed to verify request: %v", err)
			writeTokenErrorResponse(w, oidc.ErrInvalidRequest, state)
		} else {
			logger.Error("Failed to verify request: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, state)
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
			writeTokenErrorResponse(w, oidc.ErrServerError, "")
		}
		return
	}

	res, err := oidc.GenerateJWKSet(project.TokenConfig.SigningAlgorithm, project.TokenConfig.SignPublicKey)
	if err != nil {
		logger.Error("Failed to generate JWT set: %+v", err)
		writeTokenErrorResponse(w, oidc.ErrServerError, "")
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
		e, ok := err.(*oidc.Error)
		if ok {
			writeTokenErrorResponse(w, e, authReq.State)
		} else {
			logger.Error("Failed to cast to *oidc.Error, this is critical program bug: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, "")
		}
		return
	}

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
		writeTokenErrorResponse(w, oidc.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)

	authReq := oidc.NewAuthRequest(r.Form)
	if err := authReq.Validate(); err != nil {
		logger.Info("Failed to validate request: %v", err)
		e, ok := err.(*oidc.Error)
		if ok {
			writeTokenErrorResponse(w, e, authReq.State)
		} else {
			logger.Error("Failed to cast to *oidc.Error, this is critical program bug: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, "")
		}
		return
	}

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
		writeTokenErrorResponse(w, oidc.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)

	// Verify user login session code
	authReq, err := oidc.UserLoginVerify(r.Form.Get("login_verify_code"))
	if err != nil {
		logger.Info("Failed to verify user login session: %v", err)
		writeTokenErrorResponse(w, oidc.ErrRequestUnauthorized, r.Form.Get("state"))
		return
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := user.Verify(projectName, uname, passwd)
	if err != nil {
		if errors.Cause(err) == user.ErrAuthFailed {
			logger.Info("Failed to authenticate user %s: %v", uname, err)
			writeTokenErrorResponse(w, oidc.ErrRequestUnauthorized, authReq.State)
		} else {
			logger.Error("Failed to verify user: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, authReq.State)
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
		writeTokenErrorResponse(w, oidc.ErrServerError, authReq.State)
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
		writeTokenErrorResponse(w, oidc.ErrRequestUnauthorized, "")
		return
	}

	user, err := db.GetInst().UserGet(claims.Subject)
	if err != nil {
		// If token validate accepted, user absolutely exists
		logger.Error("Failed to get user: %+v", err)
		writeTokenErrorResponse(w, oidc.ErrServerError, "")
		return
	}

	res := &UserInfo{
		Subject:  claims.Subject,
		UserName: user.Name,
	}

	jwthttp.ResponseWrite(w, "UserInfoHandler", res)
}

// RevokeHandler ...
func RevokeHandler(w http.ResponseWriter, r *http.Request) {
	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		writeTokenErrorResponse(w, oidc.ErrInvalidRequestObject, "")
		return
	}

	tokenType := r.Form.Get("token_type_hint")
	if tokenType == "" {
		tokenType = "refresh_token" // default is refresh token
	}

	switch tokenType {
	case "access_token":
		// TODO(implement revocation of access token)
		writeTokenErrorResponse(w, oidc.ErrUnsupportedTokenType, r.Form.Get("state"))
	case "refresh_token":
		refreshToken := r.Form.Get("token")
		claims := &token.RefreshTokenClaims{}
		issuer := token.GetExpectIssuer(r)
		if err := token.ValidateRefreshToken(claims, refreshToken, issuer); err != nil {
			logger.Info("Failed to validate refresh token: %v", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := db.GetInst().SessionDelete(claims.SessionID); err != nil {
			e := errors.Cause(err)
			if e == model.ErrNoSuchSession || e == model.ErrSessionValidateFailed {
				logger.Info("Failed to revoke session: %v", err)
				w.WriteHeader(http.StatusOK)
			} else {
				logger.Error("Failed to revoke session: %+v", err)
				writeTokenErrorResponse(w, oidc.ErrServerError, r.Form.Get("state"))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		writeTokenErrorResponse(w, oidc.ErrUnsupportedTokenType, r.Form.Get("state"))
	}
}
