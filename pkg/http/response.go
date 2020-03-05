package http

import (
	"encoding/json"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"net/http"
)

// ResponseWrite ...
func ResponseWrite(w http.ResponseWriter, handlerName string, v interface{}) {
	w.Header().Add("Content-Type", "application/json")

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			logger.Error("Failed to encode a response for %s: %+v", handlerName, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	logger.Info("%s method successfully finished", handlerName)
}
