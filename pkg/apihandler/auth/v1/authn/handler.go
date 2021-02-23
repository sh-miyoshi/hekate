package authn

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/login"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	"github.com/sh-miyoshi/hekate/pkg/sso"
	"github.com/stretchr/stew/slice"
)

var (
	errSessionEnd = errors.New("Session end", "Session end")
)

// UserLoginHandler ...
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errMsg := "Request failed. invalid form value"
		login.WriteErrorPage(errMsg, w)
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")

	sessionID := r.Form.Get("login_session_id")

	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
			// delete session if login failed
			db.GetInst().LoginSessionDelete(projectName, sessionID)
		}

		if err = audit.GetInst().Save(projectName, time.Now(), "USER_LOGIN", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Verify user login session code
	s, err := login.VerifySession(projectName, sessionID)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to verify user login session"))
		errMsg := "Request failed. Internal server error occured."
		if errors.Contains(err, errors.ErrSessionExpired) {
			errMsg = "Request failed. The session was already expired."
		} else if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			errMsg = "Request failed. Invalid request was sent."
		}
		login.WriteErrorPage(errMsg, w)
		return
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := login.UserVerifyByPassword(projectName, uname, passwd)
	if err != nil {
		if errors.Contains(err, login.ErrAuthFailed) || errors.Contains(err, login.ErrUserLocked) {
			errors.PrintAsInfo(errors.Append(err, "Failed to authenticate user %s", uname))

			// delete old session and create new code for relogin
			if err := db.GetInst().LoginSessionDelete(projectName, sessionID); err != nil {
				errors.Print(errors.Append(err, "Failed to delete previous login session"))
				login.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}

			var authReq oidc.AuthRequest
			copier.Copy(&authReq, &s)
			authReq.State = state

			lsID, err := login.StartLoginSession(projectName, &authReq)
			if err != nil {
				errors.Print(errors.Append(err, "Failed to register login session"))
				login.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}
			login.WriteUserLoginPage(projectName, lsID, "invalid user name or password", state, w)
			err = nil // do not delete session in defer function
		} else {
			errors.Print(errors.Append(err, "Failed to verify user"))
			errMsg := "Request failed. internal server error occuerd"
			login.WriteErrorPage(errMsg, w)
		}
		return
	}

	s.UserID = usr.ID
	s.LoginDate = time.Now()

	if err = db.GetInst().LoginSessionUpdate(projectName, s); err != nil {
		errors.Print(errors.Append(err, "Failed to update login session"))
		login.WriteErrorPage("Request failed. internal server error occuerd", w)
		return
	}

	logger.Debug("Successfully verify user login by password")

	// Next Steps.
	// 1. If required MFA, return MFA page
	// 2. If required content, return consent page
	// 3. login session finished, redirect to callback URL

	// TODO MFA Page

	// Consent Page
	if ok := slice.Contains(s.Prompt, "consent"); ok {
		login.WriteConsentPage(projectName, sessionID, state, w)
		return
	}

	// Login session finished, redirect to callback URL
	req, err := redirectToCallback(w, r, projectName, s)
	if err != nil {
		if !errors.Contains(err, errSessionEnd) {
			errors.Print(err)
			login.WriteErrorPage("Request failed. internal server error occuerd", w)
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
		login.WriteErrorPage(errMsg, w)
		return
	}

	logger.Debug("Form: %v", r.Form)
	state := r.Form.Get("state")
	sessionID := r.Form.Get("login_session_id")

	var err *errors.Error
	defer func() {
		if err != nil {
			// delete session if login failed
			db.GetInst().LoginSessionDelete(projectName, sessionID)
		}
	}()

	s, err := login.VerifySession(projectName, sessionID)
	if err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to verify user login session"))
		errMsg := "Request failed. Internal server error occured."
		if errors.Contains(err, errors.ErrSessionExpired) {
			errMsg = "Request failed. The session was already expired."
		} else if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			errMsg = "Request failed. Invalid request was sent."
		}
		login.WriteErrorPage(errMsg, w)
		return
	}

	switch sel {
	case "yes":
		req, err := redirectToCallback(w, r, projectName, s)
		if err != nil {
			if !errors.Contains(err, errSessionEnd) {
				errors.Print(err)
				login.WriteErrorPage("Request failed. internal server error occuerd", w)
				return
			}
		}
		http.Redirect(w, req, req.URL.String(), http.StatusFound)
	case "no":
		err = errors.ErrConsentRequired
		errors.RedirectWithOAuthError(w, err, r.Method, s.RedirectURI, state)
	default:
		msg := fmt.Sprintf("Invalid select type %s. consent page maybe broken.", sel)
		err = errors.New("Invalid select type", msg)
		logger.Info(msg)
		login.WriteErrorPage("Request failed. internal server error occuerd", w)
	}
}

func redirectToCallback(w http.ResponseWriter, r *http.Request, projectName string, session *model.LoginSession) (*http.Request, *errors.Error) {
	state := r.Form.Get("state")
	issuer := token.GetFullIssuer(r)

	req, err := oidc.CreateLoggedInResponse(session, state, issuer)
	if err != nil {
		return nil, err
	}

	if ok := slice.Contains(session.ResponseType, "code"); !ok && len(session.ResponseType) > 0 {
		// delete session
		return req, errSessionEnd
	}

	if err := sso.SetSSOSessionToCookie(w, projectName, session.UserID, issuer); err != nil {
		return nil, errors.Append(err, "Failed to set cookie")
	}

	return req, nil
}
