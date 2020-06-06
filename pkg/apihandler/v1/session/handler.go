package sessionapi

import (
	"net/http"

	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// SessionDeleteHandler ...
//   require role: write-project
func SessionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("SessionDeleteHandler method successfully finished")
}

// SessionGetHandler ...
//   require role: read-client
func SessionGetHandler(w http.ResponseWriter, r *http.Request) {
}
