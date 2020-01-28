package oidc

import (
	"encoding/json"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
)

func writeTokenErrorResponse(w http.ResponseWriter) {
	res := ErrorResponse{
		ErrorCode: "invalid_request",
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error("Failed to encode a token error response: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
