package tokenapi

import (
	"net/http"

	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// TokenCreateHandler method create JWT token
func TokenCreateHandler(w http.ResponseWriter, r *http.Request) {

	logger.Info("TokenCreateHandler method is not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)
}
