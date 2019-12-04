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

	userID := request.ID
	if userID == "" {
		if request.Name == "" {
			logger.Info("Name or ID must be specified")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		userID, err = db.GetInst().User.GetIDByName(projectID, request.Name)
		if err != nil {
			if err == model.ErrNoSuchUser {
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				logger.Error("Failed to get user id: %+v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
	}
	logger.Debug("User ID: %s", userID)

	user, err := db.GetInst().User.Get(projectID, userID)
	if err != nil {
		if err == model.ErrNoSuchUser {
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user info: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Secret Authenticate
	switch request.AuthType {
	case "password":
		hash := util.CreateHash(request.Secret)
		if user.PasswordHash != hash {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	case "refresh":
		// Parse secret which is refresh token
		claims := &token.RefreshTokenClaims{}
		if err := token.ValidateRefreshToken(claims, request.Secret); err != nil {
			logger.Info("Failed to validate refresh token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		logger.Debug("%v", claims)

		if claims.Audience != userID {
			logger.Info("Invalid refresh token audience")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	default:
		logger.Error("Invalid Authentication Type: %s", request.AuthType)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Generate JWT Token
	accessTokenReq := token.Request{
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.AccessTokenLifeSpan),
		ProjectID:   user.ProjectID,
		UserID:      user.ID,
	}

	res := TokenResponse{}
	res.AccessExpiresIn = project.TokenConfig.AccessTokenLifeSpan
	res.AccessToken, err = token.GenerateAccessToken(accessTokenReq)
	if err != nil {
		logger.Error("Failed to get JWT token: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	refreshTokenReq := token.Request{
		ExpiredTime: time.Second * time.Duration(project.TokenConfig.RefreshTokenLifeSpan),
		ProjectID:   user.ProjectID,
		UserID:      user.ID,
	}

	res.RefreshExpiresIn = project.TokenConfig.RefreshTokenLifeSpan
	res.RefreshToken, err = token.GenerateRefreshToken(refreshTokenReq)
	if err != nil {
		logger.Error("Failed to get JWT token: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JWT Token
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		logger.Error("Failed to encode a response for JWT token create: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("TokenCreateHandler method successfully finished")
}
