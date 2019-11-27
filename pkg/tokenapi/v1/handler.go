package tokenapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"time"
)

// TokenCreateHandler method create JWT token
func TokenCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]

	// TODO(Validate project ID)

	// Get Project Info for Token Config
	project, err := db.GetInst().Project.Get(projectID)
	if err == model.ErrNoSuchProject {
		http.Error(w, "Project Not Found", http.StatusNotFound)
		return
	}

	// Parse Request
	var request TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode token create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Validate Request)

	userID, err := db.GetInst().User.GetIDByName(projectID, request.Name)
	if err != nil {
		if err == model.ErrNoSuchUser {
			http.Error(w, "User Not Found", http.StatusNotFound)
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

	// Generate JWT Token
	tokenReq := token.Request{
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		Audience:    user.ID,
	}

	res := TokenResponse{}
	res.AccessToken, err = token.Generate(tokenReq)
	if err != nil {
		logger.Error("Failed to get JWT token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JWT Token
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		logger.Error("Failed to encode a response for JWT token create: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("TokenCreateHandler method successfully finished")

}
