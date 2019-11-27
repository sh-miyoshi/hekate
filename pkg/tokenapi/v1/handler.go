package tokenapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
	"os"
)

// TokenCreateHandler method create JWT token
func TokenCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]

	// TODO(Validate project ID)

	// Parse Request
	var request TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info(projectID, "Failed to decode token create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Validate Request)

	userID, err := db.GetInst().User.GetIDByName(projectID, request.Name)
	if err != nil {
		if err == model.ErrNoSuchUser || err == os.ErrNotExist {
			http.Error(w, "User or Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user id: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	logger.Debug("User ID: %s", userID)

	user, err := db.GetInst().User.Get(projectID, userID)
	if err != nil {
		if err == model.ErrNoSuchUser {
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user info: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Secret Authenticate
	hash := util.CreateHash(request.Secret)
	switch request.AuthType {
	case "password":
		if user.PasswordHash != hash {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	default:
		logger.Error("Invalid Authentication Type: %s", request.AuthType)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Generate JWT Token)

	// Return JWT Token
	logger.Info("TokenCreateHandler method is not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}
