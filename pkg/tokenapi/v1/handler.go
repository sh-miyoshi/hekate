package tokenapi

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"time"
)

// TokenCreateHandler method create JWT token
func TokenCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// TODO(Validate project ID)

	// Get Project Info for Token Config
	project, err := db.GetInst().ProjectGet(projectName)
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

	user := &model.UserInfo{}

	if request.ID == "" {
		if request.Name == "" {
			logger.Info("Name or ID must be specified")
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user, err = db.GetInst().UserGetByName(projectName, request.Name)
	} else {
		user, err = db.GetInst().UserGet(projectName, request.ID)
	}

	if err != nil {
		if err == model.ErrNoSuchUser {
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user id: %+v", err)
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

		if claims.Audience != user.ID {
			logger.Info("Invalid refresh token audience")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Revoke previous session
		if err := db.GetInst().RevokeSession(projectName, user.ID, claims.SessionID); err != nil {
			// Not found the session
			logger.Info("Token is already revoked")
			http.Error(w, "Token Revoked", http.StatusUnauthorized)
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
		ProjectName: user.ProjectName,
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
		ProjectName: user.ProjectName,
		UserID:      user.ID,
	}

	sessionID := uuid.New().String()
	res.RefreshExpiresIn = project.TokenConfig.RefreshTokenLifeSpan
	res.RefreshToken, err = token.GenerateRefreshToken(sessionID, refreshTokenReq)
	if err != nil {
		logger.Error("Failed to get JWT token: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	session := model.Session{
		SessionID: sessionID,
		CreatedAt: now,
		ExpiresIn: res.RefreshExpiresIn,
		FromIP:    r.RemoteAddr,
	}

	if err := db.GetInst().NewSession(projectName, user.ID, session); err != nil {
		logger.Error("Failed to register refresh token session token: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jwthttp.ResponseWrite(w, "TokenCreateHandler", &res)
}
