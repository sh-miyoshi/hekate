package userv1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/secret"
)

// ChangePasswordHandler ...
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIToken(r)
	if err != nil || claims.Subject != userID {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
	}

	var req ChangePasswordRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		err = errors.New("Invalid request", "Failed to decode user change password request: %v", e)
		errors.PrintAsInfo(err)
		errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		return
	}

	if err = db.GetInst().UserChangePassword(projectName, userID, req.Password); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			logger.Info("No such user: %s", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
			} else {
				errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
				errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
			}
		} else if errors.Contains(err, secret.ErrPasswordPolicyFailed) {
			errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to change yser password"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	if err = db.GetInst().UserLogout(projectName, userID); err != nil {
		if errors.Contains(err, model.ErrUserValidateFailed) {
			logger.Info("User ID %s is invalid", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to logout"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserLogoutHandler method successfully finished")
}