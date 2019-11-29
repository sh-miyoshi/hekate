package projectapi

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
	"time"
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
	// Parse Request
	var request ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode token create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Validate Request)

	// Create Project Entry
	project := model.ProjectInfo{
		ID:        uuid.New().String(),
		Name:      request.Name,
		Enabled:   request.Enabled,
		CreatedAt: time.Now().String(),
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  request.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: request.TokenConfig.RefreshTokenLifeSpan,
		},
	}

	// Create New Project
	if err := db.GetInst().Project.Add(&project); err != nil {
		if err == model.ErrProjectAlreadyExists {
			http.Error(w, "Project Already Exists", http.StatusConflict)
		} else {
			logger.Error("Failed to create project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	w.Header().Add("Content-Type", "application/json")
	res := ProjectGetResponse{
		ID:        project.ID,
		Name:      project.Name,
		Enabled:   project.Enabled,
		CreatedAt: project.CreatedAt,
		TokenConfig: &TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
		},
	}

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		logger.Error("Failed to encode a response for JWT token create: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("ProjectCreateHandler method successfully finished")
}

// ProjectDeleteHandler ...
func ProjectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID := vars["projectID"]

	if projectID == "master" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := db.GetInst().Project.Delete(projectID); err != nil {
		if err == model.ErrNoSuchProject {
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to create project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("ProjectDeleteHandler method successfully finished")
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
