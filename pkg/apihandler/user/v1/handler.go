package userv1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/otp"
	"github.com/sh-miyoshi/hekate/pkg/secret"
)

// GetHandler ...
func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) || errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "User %s is not found", userID))
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else {
			errors.Print(errors.Append(err, "Failed to get user"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	// Return Response
	res := &GetResponse{
		ID:        user.ID,
		Name:      user.Name,
		EMail:     user.EMail,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		OPTInfo: OTPInfo{
			ID:      user.OTPInfo.ID,
			Enabled: user.OTPInfo.Enabled,
		},
	}
	jwthttp.ResponseWrite(w, "GetHandler", res)
}

// ChangePasswordHandler ...
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	var req ChangePasswordRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		err = errors.Append(errors.ErrInvalidRequest, "Failed to decode user change password request: %v", e)
		errors.PrintAsInfo(err)
		errors.WriteToHTTP(w, err, 0, "")
		return
	}

	if err = db.GetInst().UserChangePassword(projectName, userID, req.Password); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			logger.Info("No such user: %s", userID)
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				errors.WriteToHTTP(w, err, http.StatusNotFound, "")
			} else {
				errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
				errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
			}
		} else if errors.Contains(err, secret.ErrPasswordPolicyFailed) {
			errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
			errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		} else {
			errors.Print(errors.Append(err, "Failed to change user password"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("ChangePasswordHandler method successfully finished")
}

// LogoutHandler ...
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	if err = db.GetInst().UserLogout(projectName, userID); err != nil {
		if errors.Contains(err, model.ErrUserValidateFailed) {
			logger.Info("User ID %s is invalid", userID)
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else {
			errors.Print(errors.Append(err, "Failed to logout"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserLogoutHandler method successfully finished")
}

// OTPGenerateHandler ...
func OTPGenerateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	qrcode, err := otp.Register(projectName, userID, claims.UserName)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to register OTP"))
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		return
	}
	logger.Debug("Generated QR code size: %d", len(qrcode))

	// Return Response
	res := &OTPGenerateResponse{
		QRCodeImage: qrcode,
	}
	jwthttp.ResponseWrite(w, "OTPGenerateHandler", res)
}

// OTPVerifyHandler ...
func OTPVerifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	var req OTPVerifyRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		err := errors.Append(errors.ErrInvalidRequest, "Failed to decode user otp verify request: %v", e)
		errors.PrintAsInfo(err)
		errors.WriteToHTTP(w, err, 0, "")
		return
	}

	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		errors.Print(err)
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		return
	}

	if !user.OTPInfo.Enabled {
		user.OTPInfo.Enabled = true
		logger.Debug("After enabled OTP: %+v", user.OTPInfo)
		if err := db.GetInst().UserUpdate(projectName, user); err != nil {
			errors.Print(err)
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
			return
		}
	}

	if err := otp.Verify(time.Now(), user, req.UserCode); err != nil {
		if errors.Contains(err, otp.ErrVerifyFailed) {
			errors.PrintAsInfo(err)
			errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		} else {
			errors.Print(err)
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("OTPVerifyHandler method successfully finished")
}

// OTPDeleteHandler ...
func OTPDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		errors.Print(err)
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		return
	}

	// Remove OTP Settings
	user.OTPInfo.Enabled = false
	user.OTPInfo.ID = ""
	user.OTPInfo.PrivateKey = ""

	if err := db.GetInst().UserUpdate(projectName, user); err != nil {
		errors.Print(err)
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("OTPDeleteHandler method successfully finished")
}
