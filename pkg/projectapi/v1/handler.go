package projectapi

import (
	"encoding/json"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
)

// AllProjectGetHandler ...
func AllProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	projectIDs, err := db.GetInst().Project.GetList()
	if err != nil {
		logger.Error("Failed to get project list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&projectIDs); err != nil {
		logger.Error("Failed to encode a response for getting project list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("AllProjectGetHandler method successfully finished")
}

// ProjectCreateHandler ...
func ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}

// ProjectDeleteHandler ...
func ProjectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}

// ProjectGetHandler ...
func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}

// ProjectUpdateHandler ...
func ProjectUpdateHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}
