package oidc

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/client"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/user"
	"github.com/stretchr/stew/slice"
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
		ResponseModesSupported: []string{
			"query",
			"fragment",
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
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	// Get Project Info for Token Config
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchProject) {
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to get project info"))
			errors.WriteOAuthError(w, errors.ErrServerError, state)
		}
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

	if err := oidc.ClientAuth(projectName, clientID, clientSecret); err != nil {
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
		if err := client.CheckRedirectURL(projectName, clientID, r.Form.Get("redirect_uri")); err != nil {
			if errors.Contains(err, client.ErrNoRedirectURL) {
				logger.Info("Redirect URL %s is not in Allowed list", r.Form.Get("redirect_uri"))
				errors.WriteOAuthError(w, errors.ErrInvalidRequestURI, state)
			} else if errors.Contains(err, model.ErrNoSuchClient) {
				errors.PrintAsInfo(errors.Append(err, "Failed to get allowed callback urls"))
				errors.WriteOAuthError(w, errors.ErrInvalidClient, state)
			} else {
				errors.Print(errors.Append(err, "Failed to get allowed callbak urls in client"))
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
		tkn, err = oidc.ReqAuthByClientCredentials(project, clientID, r)
	case model.GrantTypePassword:
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")
		tkn, err = oidc.ReqAuthByPassword(project, uname, passwd, r)
	case model.GrantTypeRefreshToken:
		refreshToken := r.Form.Get("refresh_token")
		tkn, err = oidc.ReqAuthByRefreshToken(project, clientID, refreshToken, r)

		if err != nil && errors.Contains(err, model.ErrNoSuchSession) {
			logger.Info("Refresh token is already revoked")
			errors.WriteOAuthError(w, errors.ErrInvalidRequest, state)
			return
		}
	case model.GrantTypeAuthorizationCode:
		code := r.Form.Get("code")
		tkn, err = oidc.ReqAuthByCode(project, clientID, code, r)
	default:
		logger.Info("Unexpected grant type got: %s", gt.String())
		errors.WriteOAuthError(w, errors.ErrServerError, state)
		return
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
		if errors.Contains(err, model.ErrNoSuchProject) {
			errors.PrintAsInfo(errors.Append(err, "No such project %s", projectName))
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to get project"))
			errors.WriteOAuthError(w, errors.ErrServerError, "")
		}
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

	issuer := token.GetExpectIssuer(r)
	authHandler(w, projectName, issuer, queries)
}

// AuthPOSTHandler ...
func AuthPOSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errMsg := "Request failed. invalid form value"
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	logger.Debug("Form: %v", r.Form)
	issuer := token.GetExpectIssuer(r)
	authHandler(w, projectName, issuer, r.Form)
}

// UserLoginHandler ...
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errMsg := "Request failed. invalid form value"
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	sessionID := r.Form.Get("login_session_id")

	// delete session if login failed
	defer func() {
		if err != nil {
			db.GetInst().LoginSessionDelete(projectName, sessionID)
		}
	}()

	// Verify user login session code
	s, err := oidc.VerifySession(projectName, sessionID)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to verify user login session"))
		errMsg := "Request failed. internal server error occured."
		if errors.Contains(err, errors.ErrSessionExpired) {
			errMsg = "Request failed. the session was already expired."
		}
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := user.Verify(projectName, uname, passwd)
	if err != nil {
		if errors.Contains(err, user.ErrAuthFailed) {
			errors.PrintAsInfo(errors.Append(err, "Failed to authenticate user %s", uname))

			// delete old session and create new code for relogin
			if err := db.GetInst().LoginSessionDelete(projectName, sessionID); err != nil {
				errors.Print(errors.Append(err, "Failed to delete previous login session"))
				oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}

			authReq := &oidc.AuthRequest{
				Scope:        s.Scope,
				ResponseType: s.ResponseType,
				ClientID:     s.ClientID,
				RedirectURI:  s.RedirectURI,
				State:        state,
				Nonce:        s.Nonce,
				MaxAge:       s.MaxAge,
				ResponseMode: s.ResponseMode,
				Prompt:       s.Prompt,
			}

			lsID, err := oidc.StartLoginSession(projectName, authReq)
			if err != nil {
				errors.Print(errors.Append(err, "Failed to register login session"))
				oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}
			oidc.WriteUserLoginPage(projectName, lsID, "invalid user name or password", state, w)
			err = nil // do not delete session in defer function
		} else {
			errors.Print(errors.Append(err, "Failed to verify user"))
			errMsg := "Request failed. internal server error occuerd"
			oidc.WriteErrorPage(errMsg, w)
		}
		return
	}

	s.UserID = usr.ID
	s.LoginDate = time.Now()

	if ok := slice.Contains(s.Prompt, "consent"); ok {
		// save user id to session info
		if err = db.GetInst().LoginSessionUpdate(projectName, s); err != nil {
			errors.Print(errors.Append(err, "Failed to update login session"))
			oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
			return
		}

		// show consent page
		oidc.WriteConsentPage(projectName, sessionID, state, w)
		return
	}

	issuer := token.GetFullIssuer(r)
	req, err := createLoginRedirectInfo(s, state, issuer)
	if err != nil {
		errors.Print(err)
		oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
		return
	}

	if ok := slice.Contains(s.ResponseType, "code"); !ok {
		// delete session
		err = errors.New("Session end", "Session end")
	} else {
		// Update login session info
		if err = db.GetInst().LoginSessionUpdate(projectName, s); err != nil {
			errors.Print(errors.Append(err, "Failed to update login session"))
			oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
			return
		}
	}

	http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

// ConsentHandler ...
func ConsentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	sel := r.FormValue("select")
	logger.Info("Consent select: %s", sel)

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errMsg := "Request failed. invalid form value"
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	switch sel {
	case "yes":
		sessionID := r.Form.Get("login_session_id")
		s, err := oidc.VerifySession(projectName, sessionID)
		if err != nil {
			errors.PrintAsInfo(errors.Append(err, "Failed to verify user login session"))
			errMsg := "Request failed. internal server error occured."
			if errors.Contains(err, errors.ErrSessionExpired) {
				errMsg = "Request failed. the session was already expired."
			}
			oidc.WriteErrorPage(errMsg, w)
			return
		}

		issuer := token.GetFullIssuer(r)
		req, err := createLoginRedirectInfo(s, state, issuer)
		if err != nil {
			errors.Print(err)
			oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
			return
		}

		http.Redirect(w, req, req.URL.String(), http.StatusFound)
	case "no":
		errors.WriteOAuthError(w, errors.ErrConsentRequired, state)
	default:
		logger.Info("Invalid select type %s. consent page maybe broken.", sel)
		errors.WriteOAuthError(w, errors.ErrServerError, state)
	}
}

// UserInfoHandler ...
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	claims, err := jwthttp.ValidateAPIRequest(r)
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
			if errors.Contains(err, model.ErrNoSuchProject) || errors.Contains(err, model.ErrNoSuchSession) || errors.Contains(err, model.ErrSessionValidateFailed) {
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

func authHandler(w http.ResponseWriter, projectName string, tokenIssuer string, req url.Values) {
	authReq := oidc.NewAuthRequest(req)
	logger.Debug("Auth Request: %v", authReq)

	if err := authReq.Validate(); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to validate request"))
		if err.GetHTTPStatusCode() == 0 {
			errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		} else {
			errors.WriteOAuthError(w, err, authReq.State)
		}
		return
	}

	// Check Redirect URL
	if err := client.CheckRedirectURL(projectName, authReq.ClientID, authReq.RedirectURI); err != nil {
		if errors.Contains(err, client.ErrNoRedirectURL) {
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

	// if already logined (check by login_hint, prompt), redirect to callback
	if slice.Contains(authReq.Prompt, "none") {
		handleSSO(w, projectName, tokenIssuer, authReq)
		return // if prompt=none, never return login page
	}

	// Start session for login flow
	lsID, err := oidc.StartLoginSession(projectName, authReq)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to start login session"))
		errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		return
	}

	oidc.WriteUserLoginPage(projectName, lsID, "", authReq.State, w)
}

func handleSSO(w http.ResponseWriter, projectName string, tokenIssuer string, authReq *oidc.AuthRequest) {
	uname := authReq.LoginHint
	user, err := db.GetInst().UserGetList(projectName, &model.UserFilter{Name: uname})
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get user %s", uname))
		errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		return
	}
	if len(user) != 1 {
		logger.Info("Unexpect user num, expect 1 but got %d", len(user))
		errors.WriteOAuthError(w, errors.ErrInvalidRequest, authReq.State)
		return
	}

	userID := user[0].ID
	sessions, err := db.GetInst().SessionGetList(projectName, userID)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get session list"))
		errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
		return
	}
	if len(sessions) == 0 {
		logger.Info("No sessions, so return login_required")
		errors.WriteOAuthError(w, errors.ErrLoginRequired, authReq.State)
		return
	}

	// check max_age
	// if now > auth_time + max_age return login_required
	now := time.Now()
	for _, s := range sessions {
		if now.Before(s.LastAuthTime.Add(time.Second * time.Duration(s.AuthMaxAge))) {
			ls := &model.LoginSession{
				ResponseType: authReq.ResponseType,
				ProjectName:  projectName,
				UserID:       s.UserID,
				ClientID:     authReq.ClientID,
				Nonce:        authReq.Nonce,
				LoginDate:    s.LastAuthTime,
				RedirectURI:  authReq.RedirectURI,
				ResponseMode: authReq.ResponseMode,
			}
			req, err := createLoginRedirectInfo(ls, authReq.State, tokenIssuer)
			if err != nil {
				errors.Print(errors.Append(err, "Failed to create login redirect info"))
				errors.WriteOAuthError(w, errors.ErrServerError, authReq.State)
				return
			}
			http.Redirect(w, req, req.URL.String(), http.StatusFound)
			return
		}
	}

	logger.Info("No valid session, so return login_required")
	errors.WriteOAuthError(w, errors.ErrLoginRequired, authReq.State)
}

func createLoginRedirectInfo(session *model.LoginSession, state, tokenIssuer string) (*http.Request, *errors.Error) {
	values := url.Values{}
	if state != "" {
		values.Set("state", state)
	}

	for _, typ := range session.ResponseType {
		switch typ {
		case "code":
			code := uuid.New().String()
			session.Code = code
			values.Set("code", code)
		case "id_token":
			prj, err := db.GetInst().ProjectGet(session.ProjectName)
			if err != nil {
				return nil, errors.Append(err, "Failed to get token lifespan in project")
			}

			audiences := []string{session.UserID, session.ClientID}
			tokenReq := token.Request{
				Issuer:          tokenIssuer,
				ExpiredTime:     time.Second * time.Duration(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName:     session.ProjectName,
				UserID:          session.UserID,
				Nonce:           session.Nonce,
				EndUserAuthTime: session.LoginDate,
			}
			tkn, err := token.GenerateIDToken(audiences, tokenReq)
			if err != nil {
				return nil, errors.Append(err, "Failed to generate id token")
			}
			values.Set("id_token", tkn)
		case "token":
			prj, err := db.GetInst().ProjectGet(session.ProjectName)
			if err != nil {
				return nil, errors.Append(err, "Failed to get token lifespan in project")
			}

			audiences := []string{session.UserID, session.ClientID}
			tokenReq := token.Request{
				Issuer:      tokenIssuer,
				ExpiredTime: time.Second * time.Duration(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName: session.ProjectName,
				UserID:      session.UserID,
			}
			tkn, err := token.GenerateAccessToken(audiences, tokenReq)
			if err != nil {
				return nil, errors.Append(err, "Failed to generate access token")
			}
			values.Set("access_token", tkn)
		default:
			return nil, errors.New("Unknown response type", "Unknown response type %s", typ)
		}
	}

	req, err := http.NewRequest("GET", session.RedirectURI, nil)
	if err != nil {
		return nil, errors.New("Internal server error", "Failed to create response: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if session.ResponseMode == "query" {
		req.URL.RawQuery = values.Encode()
	} else if session.ResponseMode == "fragment" {
		req.URL.Fragment = values.Encode()
	} else {
		return nil, errors.New("Internal server error", "Invalid response mode %s is specified", session.ResponseMode)
	}

	return req, nil
}
