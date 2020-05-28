package oidc

import (
	"encoding/json"
	"net/http"

	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
)

func writeErrorResponse(w http.ResponseWriter, err *oidc.Error, state string) {
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
