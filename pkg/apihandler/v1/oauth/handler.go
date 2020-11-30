package oauth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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
	"github.com/sh-miyoshi/hekate/pkg/util"
	"github.com/stretchr/stew/slice"
)

const (
	userCodeLength      = 8
	tokenReqIntervalSec = 5
)

// DeviceRegisterHandler ...
func DeviceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "DEVICE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}

	logger.Debug("Form: %v", r.Form)

	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")

	if clientID == "" {
		// maybe basic authentication
		i, s, ok := r.BasicAuth()
		if !ok {
			logger.Info("Failed to get client ID from request, Request header: %v", r.Header)
			errors.WriteOAuthError(w, errors.ErrInvalidClient, "")
			return
		}
		clientID = i
		clientSecret = s
	}

	if err = oidc.ClientAuth(projectName, clientID, clientSecret); err != nil {
		if errors.Contains(err, errors.ErrInvalidClient) {
			errors.PrintAsInfo(errors.Append(err, "Failed to authenticate client %s", clientID))
			errors.WriteOAuthError(w, errors.ErrInvalidClient, "")
		} else {
			errors.Print(errors.Append(err, "Failed to authenticate client"))
			errors.WriteOAuthError(w, errors.ErrServerError, "")
		}
		return
	}

	scope := r.Form.Get("scope")
	if !slice.Contains(config.Get().SupportedScope, scope) {
		errors.PrintAsInfo(errors.New("Invalid scope", "Invalid scope request: %s", scope))
		errors.WriteOAuthError(w, errors.ErrRequestNotSupported, "")
		return
	}

	url := config.GetServerAddr(r) + "/resource/project/" + projectName + "/devicecomplete"
	authReq := &oidc.AuthRequest{
		Scope:        scope,
		ClientID:     clientID,
		RedirectURI:  url,
		ResponseMode: "query",
	}

	lsID, err := login.StartLoginSession(projectName, authReq)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to start device login session"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	deviceCode := uuid.New().String()
	userCode := util.RandomString(userCodeLength, util.CharTypeUpper)
	expires := config.Get().LoginSessionExpiresIn

	ent := &model.Device{
		DeviceCode:     deviceCode,
		UserCode:       userCode,
		ProjectName:    projectName,
		CreatedAt:      time.Now(),
		ExpiresIn:      int64(expires),
		LoginSessionID: lsID,
	}

	if err := db.GetInst().DeviceAdd(projectName, ent); err != nil {
		errors.Print(errors.Append(err, "Failed to add device"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	url = config.GetServerAddr(r) + "/resource/project/" + projectName + "/devicelogin"
	res := DeviceAuthorizationResponse{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationURI: url,
		ExpiresIn:       int(expires),
		Interval:        tokenReqIntervalSec,
	}

	logger.Debug("Device Authorization Response: %v", res)
	jwthttp.ResponseWrite(w, "DeviceHandler", &res)
}

// DeviceLoginPageHandler return html page for input user code
func DeviceLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// set error if exists
	queries := r.URL.Query()
	err := queries.Get("error")

	login.WriteDeviceLoginPage(projectName, err, w)
}

// DeviceUserCodeVerifyHandler ...
func DeviceUserCodeVerifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "DEVICE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	if err := r.ParseForm(); err != nil {
		logger.Info("Failed to parse form: %v", err)
		errors.WriteOAuthError(w, errors.ErrInvalidRequestObject, "")
		return
	}
	userCode := r.Form.Get("code")
	devices, err := db.GetInst().DeviceGetList(projectName, &model.DeviceFilter{UserCode: userCode})
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get device"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}
	if len(devices) == 0 {
		logger.Info("No valid device for user code: %s", userCode)
		login.WriteDeviceLoginPage(projectName, "The code is invalid", w)
		return
	}

	// ok to verify user code, next is user authentication
	login.WriteUserLoginPage(projectName, devices[0].LoginSessionID, "", "", w)
}
