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

	authReq := &oidc.AuthRequest{
		Scope:       scope,
		ClientID:    clientID,
		RedirectURI: "http://localhost:8080/device/complete", // TODO
	}

	lsID, err := login.StartLoginSession(projectName, authReq)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to start device login session"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	length := 8 // TODO use const value
	deviceCode := uuid.New().String()
	userCode := util.RandomString(length, util.CharTypeUpper)
	expires := int64(600) // TODO set corrent value

	ent := &model.Device{
		DeviceCode:     deviceCode,
		UserCode:       userCode,
		ProjectName:    projectName,
		CreatedAt:      time.Now(),
		ExpiresIn:      expires,
		LoginSessionID: lsID,
	}

	if err := db.GetInst().DeviceAdd(projectName, ent); err != nil {
		errors.Print(errors.Append(err, "Failed to add device"))
		errors.WriteOAuthError(w, errors.ErrServerError, "")
		return
	}

	res := DeviceAuthorizationResponse{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationURI: "http://localhost:8080/project/" + projectName + "/devicelogin", // TODO set correct value
		ExpiresIn:       int(expires),
		Interval:        5, // TODO set correct value
	}

	logger.Debug("Device Authorization Response: %v", res)
	jwthttp.ResponseWrite(w, "DeviceHandler", &res)
}

// DeviceLoginPageHandler return html page for input user code
func DeviceLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// TODO get error from query
	login.WriteDeviceLoginPage(projectName, "", w)
}

// DeviceUserCodeVerifyHandler ...
func DeviceUserCodeVerifyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	// get user_code, projectName
	// if ok return login page
	w.Write([]byte("ok"))
}
