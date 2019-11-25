package tokenapi

import (
	"net/http"

	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// TokenCreateHandler method create JWT token
func TokenCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	// Validate Request
	// Password Authenticate
	// Return JWT Token
	logger.Info("TokenCreateHandler method is not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}
