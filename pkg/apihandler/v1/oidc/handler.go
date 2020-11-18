package oidc

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/login"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/authn"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/sso"
	"github.com/stretchr/stew/slice"
)

// ConfigGetHandler method return a configuration of OpenID Connect
func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	issuer := token.GetFullIssuer(r)
	logger.Debug("Issuer: %s", issuer)

	prj, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get project info"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}
	grantTypes := []string{}
	for _, t := range prj.AllowGrantTypes {
		grantTypes = append(grantTypes, string(t))
	}

	cfg := config.Get()
	res := Config{
		Issuer:                 issuer,
		AuthorizationEndpoint:  issuer + "/openid-connect/auth",
		TokenEndpoint:          issuer + "/openid-connect/token",
		UserinfoEndpoint:       issuer + "/openid-connect/userinfo",
		JwksURI:                issuer + "/openid-connect/certs",
		ScopesSupported:        cfg.SupportedScope,
		ResponseTypesSupported: cfg.SupportedResponseType,
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
		ResponseModesSupported: []string{
			"query",
			"fragment",
		},
		GrantTypesSupported: grantTypes,
		TokenEndpointAuthMethodsSupported: []string{
			"client_secret_basic",
			"client_secret_post",
		},
	}

	jwthttp.ResponseWrite(w, "ConfigGetHandler", &res)
}

// TokenHandler ...
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "TOKEN", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	// Get Project Info for Token Config
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get project info"))
		errors.WriteOAuthError(w, errors.ErrServerError, state)
		return
	}

	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	if clientID == "" {
		// maybe basic authentication
		i, s, ok := r.BasicAuth()
		if !ok {
			logger.Info("Failed to get client ID from request, Request header: %v", r.Header)
			errors.WriteOAuthError(w, errors.ErrInvalidClient, state)
			return
		}
		clientID = i
		clientSecret = s
	}

	if err = oidc.ClientAuth(projectName, clientID, clientSecret); err != nil {
		if errors.Contains(err, errors.ErrInvalidClient) {
			errors.PrintAsInfo(errors.Append(err, "Failed to authenticate client %s", clientID))
			errors.WriteOAuthError(w, errors.ErrInvalidClient, state)
		} else {
			errors.Print(errors.Append(err, "Failed to authenticate client"))
			errors.WriteOAuthError(w, errors.ErrServerError, state)
		}
		return
	}

	var tkn *oidc.TokenResponse

	if r.Form.Get("redirect_uri") != "" {
		// existence of client is already checked in oidc.ClientAuth
		if err = oidc.CheckRedirectURL(projectName, clientID, r.Form.Get("redirect_uri")); err != nil {
			if errors.Contains(err, oidc.ErrNoRedirectURL) {
				logger.Info("Redirect URL %s is not in Allowed list", r.Form.Get("redirect_uri"))
				errors.WriteOAuthError(w, errors.ErrInvalidRequestURI, state)
			} else {
				errors.Print(errors.Append(err, "Failed to get allowed callback urls in client"))
				errors.WriteOAuthError(w, errors.ErrServerError, state)
			}
			return
		}
	}

	// Authetication
	gtStr := r.Form.Get("grant_type")
	gt, err := model.GetGrantType(gtStr)
	if err != nil {
		logger.Info("No such Grant Type: %s", gtStr)
		errors.WriteOAuthError(w, errors.ErrInvalidGrant, state)
		return
	}
	if ok := slice.Contains(project.AllowGrantTypes, gt); !ok {
		logger.Info("Grant Type %s is not in allowed list %v", gtStr, project.AllowGrantTypes)
		errors.WriteOAuthError(w, errors.ErrUnsupportedGrantType, state)
	}

	switch gt {
	case model.GrantTypeClientCredentials:
		tkn, err = authn.ReqAuthByClientCredentials(project, clientID, r)
	case model.GrantTypePassword:
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")
		tkn, err = authn.ReqAuthByPassword(project, uname, passwd, r)
	case model.GrantTypeRefreshToken:
		refreshToken := r.Form.Get("refresh_token")
		tkn, err = authn.ReqAuthByRefreshToken(project, clientID, refreshToken, r)

		if err != nil && errors.Contains(err, model.ErrNoSuchSession) {
			logger.Info("Refresh token is already revoked")
			errors.WriteOAuthError(w, errors.ErrInvalidGrant, state)
			return
		}
	case model.GrantTypeAuthorizationCode:
		code := r.Form.Get("code")
		codeVerifier := r.Form.Get("code_verifier")
		tkn, err = authn.ReqAuthByCode(project, clientID, code, codeVerifier, r)
	case model.GrantTypeDevice:
		deviceCode := r.Form.Get("device_code")
		tkn, err = authn.ReqAuthByDeviceCode(project, clientID, deviceCode, r)
	}

	if err != nil {
		if err.GetHTTPStatusCode() != 0 {
			errors.PrintAsInfo(errors.Append(err, "Failed to verify request"))
			errors.WriteOAuthError(w, err, state)
		} else {
			errors.Print(errors.Append(err, "Failed to verify request"))
			errors.WriteOAuthError(w, errors.ErrServerError, state)
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
		errors.Print(errors.Append(err, "Failed to get project"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	res, err := oidc.GenerateJWKSet(project.TokenConfig.SigningAlgorithm, project.TokenConfig.SignPublicKey)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to generate JWT set"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Pragma", "no-cache")
	jwthttp.ResponseWrite(w, "CertsHandler", res)
}

// AuthGETHandler ...
func AuthGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Query
	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	authHandler(w, r, projectName, queries)
}

// AuthPOSTHandler ...
func AuthPOSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)
	authHandler(w, r, projectName, r.Form)
}

// UserInfoHandler ...
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to validate header"))
		errors.WriteOAuthError(w, errors.ErrInvalidRequest, "")
		return
	}

	user, err := db.GetInst().UserGet(projectName, claims.Subject)
	if err != nil {
		// If token validate accepted, user absolutely exists
		errors.Print(errors.Append(err, "Failed to get user"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	res := &UserInfo{
		Subject:  claims.Subject,
		UserName: user.Name,
	}

	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Pragma", "no-cache")
	jwthttp.ResponseWrite(w, "UserInfoHandler", res)
}

// RevokeHandler ...
func RevokeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}

	tokenType := r.Form.Get("token_type_hint")
	if tokenType == "" {
		tokenType = "refresh_token" // default is refresh token
	}

	switch tokenType {
	case "access_token":
		errors.WriteOAuthError(w, errors.ErrUnsupportedTokenType, r.Form.Get("state"))
	case "refresh_token":
		refreshToken := r.Form.Get("token")
		claims := &token.RefreshTokenClaims{}
		issuer := token.GetExpectIssuer(r)
		if err := token.ValidateRefreshToken(claims, refreshToken, issuer); err != nil {
			errors.PrintAsInfo(errors.Append(err, "Failed to validate refresh token"))
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := db.GetInst().SessionDelete(projectName, claims.SessionID); err != nil {
			if errors.Contains(err, model.ErrNoSuchSession) || errors.Contains(err, model.ErrSessionValidateFailed) {
				errors.PrintAsInfo(errors.Append(err, "Failed to revoke session"))
				w.WriteHeader(http.StatusOK)
			} else {
				errors.Print(errors.Append(err, "Failed to revoke session"))
				errors.WriteOAuthError(w, errors.ErrServerError, r.Form.Get("state"))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		errors.WriteOAuthError(w, errors.ErrUnsupportedTokenType, r.Form.Get("state"))
	}
}

func authHandler(w http.ResponseWriter, r *http.Request, projectName string, req url.Values) {
	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "AUTHORIZATION_CODE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	tokenIssuer := token.GetExpectIssuer(r)
	authReq := oidc.NewAuthRequest(req)
	logger.Debug("Auth Request: %v", authReq)

	// Check Redirect URL
	if err = oidc.CheckRedirectURL(projectName, authReq.ClientID, authReq.RedirectURI); err != nil {
		if errors.Contains(err, oidc.ErrNoRedirectURL) {
			errors.PrintAsInfo(errors.Append(err, "Redirect URL %s is not in Allowed list", authReq.RedirectURI))
			errors.WriteOAuthError(w, errors.ErrInvalidRequestURI, authReq.State)
		} else if errors.Contains(err, model.ErrNoSuchClient) {
			errors.PrintAsInfo(errors.Append(err, "Failed to get allowed callback urls: No such client %s", authReq.ClientID))
			errors.WriteOAuthError(w, errors.ErrInvalidClient, authReq.State)
		} else {
			errors.Print(errors.Append(err, "Failed to get allowed callback urls in client"))
			errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		}
		return
	}

	if err = authReq.Validate(); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to validate request"))
		if err.GetHTTPStatusCode() == 0 {
			errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		} else {
			errors.RedirectWithOAuthError(w, err, r.Method, authReq.RedirectURI, authReq.State)
		}
		return
	}

	// if prompt contains login or select_account or consent
	//   create login_session and return login page
	// else
	//   get user id from id_token_hint or cookie
	//   find sessions
	//   if ok
	//     return success response(token, code, ...)
	//   else if prompt is none
	//     return logiin_required
	//   else
	//     create login_session and return login page

	if slice.Contains(authReq.Prompt, "login") || slice.Contains(authReq.Prompt, "select_account") || slice.Contains(authReq.Prompt, "consent") {
		// Start session for login flow
		lsID, err := login.StartLoginSession(projectName, authReq)
		if err != nil {
			errors.Print(errors.Append(err, "Failed to start login session"))
			errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
			return
		}

		login.WriteUserLoginPage(projectName, lsID, "", authReq.State, w)
		return
	}

	// get user id from id_token_hint or cookie
	userID := ""
	if authReq.IDTokenHint != "" {
		var claims token.IDTokenClaims
		if err = token.ValidateIDToken(&claims, authReq.IDTokenHint, projectName, tokenIssuer); err != nil {
			errors.PrintAsInfo(errors.Append(err, "Failed to validate id_token_hint"))
			errors.RedirectWithOAuthError(w, errors.ErrInvalidRequest, r.Method, authReq.RedirectURI, authReq.State)
			return
		}
		userID = claims.Subject
	} else {
		cookie, err := r.Cookie("HEKATE_LOGIN_SESSION")
		if err != nil {
			logger.Debug("Failed to get user id from cookie: %v", err)
		} else {
			userID, err = sso.GetLoginUserIDFromSSOSessionCookie(cookie, projectName)
		}
	}

	if userID != "" {
		req, err := sso.Handle(r.Method, projectName, userID, tokenIssuer, authReq)
		if err == nil {
			http.Redirect(w, req, req.URL.String(), http.StatusFound)
			return
		} else if !errors.Contains(err, errors.ErrLoginRequired) {
			// Internal Server Error
			errors.Print(errors.Append(err, "Failed to handler SSO"))
			errors.RedirectWithOAuthError(w, errors.ErrServerError, r.Method, authReq.RedirectURI, authReq.State)
			return
		}
	}

	if slice.Contains(authReq.Prompt, "none") {
		logger.Info("request is prompt=none, but no valid sessions")
		errors.RedirectWithOAuthError(w, errors.ErrLoginRequired, r.Method, authReq.RedirectURI, authReq.State)
		return // if prompt=none, never return login page
	}

	// Start session for login flow
	lsID, err := login.StartLoginSession(projectName, authReq)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to start login session"))
		errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		return
	}

	// Return login page
	login.WriteUserLoginPage(projectName, lsID, "", authReq.State, w)
}
