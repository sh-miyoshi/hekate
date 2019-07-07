package tokenapi

import (
	"encoding/json"
	"net/http"

	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
)

// CreateTokenHandler create a token
func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("call CreateTokenHandler method with Body: %v", r.Body)

	// Parse request body
	var req TokenCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info("Failed to decode Create Token params: %v", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Authenticate request user
	if err := db.GetInst().Authenticate(req.ID, req.Password); err != nil {
		logger.Info("Failed to authenticate user: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create jwt token
	resToken, err := token.Generate()
	if err != nil {
		logger.Info("Failed to generate token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	res := TokenCreateResponse{
		Token: resToken,
	}

	resRaw, err := json.Marshal(res)
	if err != nil {
		logger.Error("Failed to marshal hobby %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resRaw)
	logger.Info("Successfully finished CreateTokenHandler")
}
