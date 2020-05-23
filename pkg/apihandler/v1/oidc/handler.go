package oidc

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/client"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
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
		writeTokenErrorResponse(w, oidc.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	// Get Project Info for Token Config
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project info: %+v", err)
			writeTokenErrorResponse(w, oidc.ErrServerError, state)
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
			writeTokenErrorResponse(w, oidc.ErrInvalidClient, state)
			return
		}
		clientID = i
		clientSecret = s
	}

	if err := oidc.ClientAuth(projectName, clientID, clientSecret); err != nil {
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

	if r.Form.Get("redirect_uri") != "" {
		if err := client.CheckRedirectURL(projectName, clientID, r.Form.Get("redirect_uri")); err != nil {
			errMsg := ""
			if errors.Cause(err) == client.ErrNoRedirectURL {
				logger.Info("Redirect URL %s is not in Allowed list", r.Form.Get("redirect_uri"))
				errMsg = "Request failed. the redirect url is not allowed"
			} else if errors.Cause(err) == model.ErrNoSuchClient {
				logger.Info("Failed to get allowed callback urls: No such client %s", clientID)
				errMsg = "Request faild. no such client"
			} else {
				logger.Error("Failed to get allowed callback urls in client: %+v", err)
				errMsg = "Request faild. internal server error occured"
			}
			oidc.WriteErrorPage(errMsg, w)
			return
		}
	}

	// Authetication
	gtStr := r.Form.Get("grant_type")
	gt, err := model.GetGrantType(gtStr)
	if err != nil {
		logger.Info("No such Grant Type: %s", gtStr)
		writeTokenErrorResponse(w, oidc.ErrInvalidGrant, state)
		return
	}
	if ok := slice.Contains(project.AllowGrantTypes, gt); !ok {
		logger.Info("Grant Type %s is not in allowed list %v", gtStr, project.AllowGrantTypes)
		writeTokenErrorResponse(w, oidc.ErrUnsupportedGrantType, state)
	}

	switch gt {
	case model.GrantTypeClientCredentials:
		tkn, err = oidc.ReqAuthByRClientCredentials(project, clientID, r)
	case model.GrantTypePassword:
		uname := r.Form.Get("username")
		passwd := r.Form.Get("password")
		tkn, err = oidc.ReqAuthByPassword(project, uname, passwd, r)
	case model.GrantTypeRefreshToken:
		refreshToken := r.Form.Get("refresh_token")
		tkn, err = oidc.ReqAuthByRefreshToken(project, clientID, refreshToken, r)

		if err != nil && errors.Cause(err) == model.ErrNoSuchSession {
			logger.Info("Refresh token is already revoked")
			writeTokenErrorResponse(w, oidc.ErrInvalidRequest, state)
			return
		}
	case model.GrantTypeAuthorizationCode:
		codeID := r.Form.Get("code")
		tkn, err = oidc.ReqAuthByCode(project, clientID, codeID, r)
	default:
		logger.Info("Unexpected grant type got: %s", gt.String())
		writeTokenErrorResponse(w, oidc.ErrServerError, state)
		return
	}

	if err != nil {
		e, ok := errors.Cause(err).(*oidc.Error)
		if ok {
			logger.Info("Failed to verify request: %v", err)
			writeTokenErrorResponse(w, e, state)
		} else {
			logger.Error("Failed to verify request: %v", err)
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

	authHandler(w, projectName, queries)
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
	authHandler(w, projectName, r.Form)
}

// UserLoginHandler ...
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
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
	state := r.Form.Get("state")

	// Verify user login session code
	info, err := oidc.UserLoginVerify(r.Form.Get("login_verify_code"))
	if err != nil {
		logger.Info("Failed to verify user login session: %v", err)
		errMsg := "Request failed. failed to verify login code"
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	// TODO(consider state)

	// TODO(set prompt value)
	authReq := &oidc.AuthRequest{
		Scope:        info.Scope,
		ResponseType: info.ResponseType,
		ClientID:     info.ClientID,
		RedirectURI:  info.RedirectURI,
		State:        state,
		Nonce:        info.Nonce,
		MaxAge:       info.MaxAge,
		ResponseMode: info.ResponseMode,
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := user.Verify(projectName, uname, passwd)
	if err != nil {
		if errors.Cause(err) == user.ErrAuthFailed {
			logger.Info("Failed to authenticate user %s: %v", uname, err)
			// create new code for relogin
			code, err := oidc.RegisterUserLoginSession(projectName, authReq)
			if err != nil {
				logger.Error("Failed to register login session %+v", err)
				oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}
			oidc.WriteUserLoginPage(projectName, code, "invalid user name or password", state, w)
		} else {
			logger.Error("Failed to verify user: %+v", err)
			errMsg := "Request failed. internal server error occuerd"
			oidc.WriteErrorPage(errMsg, w)
		}
		return
	}

	values := url.Values{}
	if state != "" {
		values.Set("state", state)
	}

	for _, typ := range info.ResponseType {
		switch typ {
		case "code":
			code, _ := oidc.GenerateAuthCode(usr.ID, *authReq)
			values.Set("code", code)
		case "id_token":
			prj, err := db.GetInst().ProjectGet(projectName)
			if err != nil {
				logger.Error("Failed to get token lifespan in project: %+v", err)
				errMsg := "Request failed. internal server error occuerd"
				oidc.WriteErrorPage(errMsg, w)
				return
			}

			audiences := []string{usr.ID, authReq.ClientID}
			tokenReq := token.Request{
				Issuer:      token.GetFullIssuer(r),
				ExpiredTime: time.Second * time.Duration(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName: projectName,
				UserID:      usr.ID,
				Nonce:       info.Nonce,
				MaxAge:      &info.MaxAge,
			}
			tkn, err := token.GenerateIDToken(audiences, tokenReq)
			if err != nil {
				logger.Error("Failed to generate id token: %+v", err)
				errMsg := "Request failed. internal server error occuerd"
				oidc.WriteErrorPage(errMsg, w)
				return
			}
			values.Set("id_token", tkn)
		case "token":
			prj, err := db.GetInst().ProjectGet(projectName)
			if err != nil {
				logger.Error("Failed to get token lifespan in project: %+v", err)
				errMsg := "Request failed. internal server error occuerd"
				oidc.WriteErrorPage(errMsg, w)
				return
			}

			audiences := []string{usr.ID, authReq.ClientID}
			tokenReq := token.Request{
				Issuer:      token.GetFullIssuer(r),
				ExpiredTime: time.Second * time.Duration(prj.TokenConfig.AccessTokenLifeSpan),
				ProjectName: projectName,
				UserID:      usr.ID,
			}
			tkn, err := token.GenerateAccessToken(audiences, tokenReq)
			if err != nil {
				logger.Error("Failed to generate access token: %+v", err)
				errMsg := "Request failed. internal server error occuerd"
				oidc.WriteErrorPage(errMsg, w)
				return
			}
			values.Set("access_token", tkn)
		default:
			logger.Error("Unknown response type: %s", typ)
			errMsg := "Request failed. internal server error occuerd"
			oidc.WriteErrorPage(errMsg, w)
			return
		}
	}

	req, err := http.NewRequest("GET", info.RedirectURI, nil)
	if err != nil {
		logger.Error("Failed to create response: %v", err)
		errMsg := "Request failed. internal server error occuerd"
		oidc.WriteErrorPage(errMsg, w)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if info.ResponseMode == "query" {
		req.URL.RawQuery = values.Encode()
	} else if info.ResponseMode == "fragment" {
		req.URL.Fragment = values.Encode()
	} else {
		logger.Error("Invalid response mode %s is specified.", info.ResponseMode)
		errMsg := "Request failed. internal server error occuerd"
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

// ConsentHandler ...
func ConsentHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(implement this)
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

func authHandler(w http.ResponseWriter, projectName string, req url.Values) {
	authReq := oidc.NewAuthRequest(req)
	if err := authReq.Validate(); err != nil {
		logger.Info("Failed to validate request: %v", err)
		errMsg := "Request failed. " + err.Error()
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	// Check Redirect URL
	if err := client.CheckRedirectURL(projectName, authReq.ClientID, authReq.RedirectURI); err != nil {
		errMsg := ""
		if errors.Cause(err) == client.ErrNoRedirectURL {
			logger.Info("Redirect URL %s is not in Allowed list", authReq.RedirectURI)
			errMsg = "Request failed. the redirect url is not allowed"
		} else if errors.Cause(err) == model.ErrNoSuchClient {
			logger.Info("Failed to get allowed callback urls: No such client %s", authReq.ClientID)
			errMsg = "Request faild. no such client"
		} else {
			logger.Error("Failed to get allowed callback urls in client: %+v", err)
			errMsg = "Request faild. internal server error occured"
		}
		oidc.WriteErrorPage(errMsg, w)
		return
	}

	// return end user auth prompt
	code, err := oidc.RegisterUserLoginSession(projectName, authReq)
	if err != nil {
		logger.Error("Failed to register login session %+v", err)
		oidc.WriteErrorPage("Request failed. internal server error occuerd", w)
		return
	}

	oidc.WriteUserLoginPage(projectName, code, "", authReq.State, w)
}
