package oidc

import (
	"encoding/json"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/oidc"
	"net/http"
)

func writeTokenErrorResponse(w http.ResponseWriter, err *oidc.Error, state string) {
	res := ErrorResponse{
		ErrorCode:   err.Name,
		Description: err.Description,
		State:       state,
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(err.Code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a token error response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
