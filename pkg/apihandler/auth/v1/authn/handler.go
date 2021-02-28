package authn

import (
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
	"github.com/sh-miyoshi/hekate/pkg/otp"
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
		errors.WriteOAuthError(w, errors.ErrServerError, "")
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
		err = errors.ErrServerError
		if errors.Contains(err, errors.ErrSessionExpired) {
			err = errors.ErrSessionExpired
		} else if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			err = errors.ErrInvalidRequest
		}
		errors.WriteOAuthError(w, err, state)
		return
	}

	// Verify user
	uname := r.Form.Get("username")
	passwd := r.Form.Get("password")
	usr, err := login.UserVerifyByPassword(projectName, uname, passwd)
	if err != nil {
		if errors.Contains(err, login.ErrAuthFailed) || errors.Contains(err, login.ErrUserLocked) {
			errors.PrintAsInfo(errors.Append(err, "Failed to authenticate user %s", uname))

			lsID, err := renewSession(projectName, s, state)
			if err != nil {
				errors.Print(err)
				errors.WriteOAuthError(w, errors.ErrServerError, state)
				return
			}

			login.WriteUserLoginPage(projectName, lsID, "invalid user name or password", state, w)
			err = nil // do not delete session in defer function
		} else {
			errors.Print(errors.Append(err, "Failed to verify user"))
			errors.WriteOAuthError(w, errors.ErrServerError, state)
		}
		return
	}

	s.UserID = usr.ID
	s.LoginDate = time.Now()

	if err = db.GetInst().LoginSessionUpdate(projectName, s); err != nil {
		errors.Print(errors.Append(err, "Failed to update login session"))
		errors.WriteOAuthError(w, errors.ErrServerError, state)
		return
	}

	logger.Debug("Successfully verify user login by password")

	// Next Steps.
	// 1. If required MFA, return MFA page
	// 2. If required content, return consent page
	// 3. login session finished, redirect to callback URL

	// OTP Verify Page
	if usr.OTPInfo.Enabled {
		login.WriteOTPVerifyPage(projectName, sessionID, state, w)
		return
	}

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
			errors.WriteOAuthError(w, errors.ErrServerError, state)
			return
		}
	}

	http.Redirect(w, req, req.URL.String(), http.StatusFound)
}

// OTPVerifyHandler ...
func OTPVerifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	userCode := r.FormValue("code")
	logger.Debug("User Code: %s", userCode)

	// Get data form Form
	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequest, "")
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
		err = errors.ErrServerError
		if errors.Contains(err, errors.ErrSessionExpired) {
			err = errors.ErrSessionExpired
		} else if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			err = errors.ErrInvalidRequest
		}
		errors.WriteOAuthError(w, err, state)
		return
	}

	user, err := db.GetInst().UserGet(projectName, s.UserID)
	if err != nil {
		err = errors.ErrServerError
		errors.WriteOAuthError(w, err, state)
		return
	}

	if err := otp.Verify(time.Now(), user, userCode); err != nil {
		if errors.Contains(err, otp.ErrVerifyFailed) {
			errors.PrintAsInfo(err)

			lsID, err := renewSession(projectName, s, state)
			if err != nil {
				errors.Print(err)
				errors.WriteOAuthError(w, errors.ErrServerError, state)
				return
			}
			// write OTP verify page again
			login.WriteOTPVerifyPage(projectName, lsID, state, w)
			return
		}
		errors.Print(err)
		errors.WriteOAuthError(w, errors.ErrServerError, state)
		return
	}

	// Next Steps.
	// 1. If required content, return consent page
	// 2. login session finished, redirect to callback URL

	// Consent Page
	if ok := slice.Contains(s.Prompt, "consent"); ok {
		login.WriteConsentPage(projectName, sessionID, state, w)
		return
	}

	// Login Success
	req, err := redirectToCallback(w, r, projectName, s)
	if err != nil {
		if !errors.Contains(err, errSessionEnd) {
			errors.Print(err)
			errors.WriteOAuthError(w, errors.ErrServerError, state)
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
		errors.WriteOAuthError(w, errors.ErrInvalidRequest, "")
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
		err = errors.ErrServerError
		if errors.Contains(err, errors.ErrSessionExpired) {
			err = errors.ErrSessionExpired
		} else if errors.Contains(err, model.ErrLoginSessionValidationFailed) {
			err = errors.ErrInvalidRequest
		}
		errors.WriteOAuthError(w, err, state)
		return
	}

	switch sel {
	case "yes":
		req, err := redirectToCallback(w, r, projectName, s)
		if err != nil {
			if !errors.Contains(err, errSessionEnd) {
				errors.Print(err)
				errors.WriteOAuthError(w, errors.ErrServerError, state)
				return
			}
		}
		http.Redirect(w, req, req.URL.String(), http.StatusFound)
	case "no":
		err = errors.ErrConsentRequired
		errors.RedirectWithOAuthError(w, err, r.Method, s.RedirectURI, state)
	default:
		err = errors.ErrServerError
		logger.Error("Invalid select type %s. consent page maybe broken.", sel)
		errors.WriteOAuthError(w, err, state)
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

func renewSession(projectName string, oldSession *model.LoginSession, state string) (string, *errors.Error) {
	// delete old session and create new code for relogin
	if err := db.GetInst().LoginSessionDelete(projectName, oldSession.SessionID); err != nil {
		return "", errors.Append(err, "Failed to delete previous login session")
	}

	var authReq oidc.AuthRequest
	copier.Copy(&authReq, &oldSession)
	authReq.State = state

	sid, err := login.StartLoginSession(projectName, &authReq)
	if err != nil {
		return "", errors.Append(err, "Failed to start new session")
	}
	return sid, nil
}
