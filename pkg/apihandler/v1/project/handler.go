package projectapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// AllProjectGetHandler ...
//   require role: cluster-read
func AllProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r, role.ResCluster, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	projects, err := db.GetInst().ProjectGetList()
	if err != nil {
		logger.Error("Failed to get project list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []ProjectGetResponse{}
	for _, prj := range projects {
		res = append(res, ProjectGetResponse{
			Name:      prj.Name,
			CreatedAt: prj.CreatedAt,
			TokenConfig: TokenConfig{
				AccessTokenLifeSpan:  prj.TokenConfig.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: prj.TokenConfig.RefreshTokenLifeSpan,
				SigningAlgorithm:     prj.TokenConfig.SigningAlgorithm,
			},
		})
	}
	logger.Debug("Project List: %v", res)

	jwthttp.ResponseWrite(w, "AllProjectGetHandler", &res)
}

// ProjectCreateHandler ...
//   require role: cluster-write
func ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r, role.ResCluster, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode project create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create Project Entry
	project := model.ProjectInfo{
		Name:         request.Name,
		CreatedAt:    time.Now(),
		PermitDelete: true,
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  request.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: request.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     request.TokenConfig.SigningAlgorithm,
		},
	}

	// Create New Project
	if err := db.GetInst().ProjectAdd(&project); err != nil {
		if errors.Cause(err) == model.ErrProjectAlreadyExists {
			logger.Info("Project %s is already exists", request.Name)
			http.Error(w, "Project Already Exists", http.StatusConflict)
		} else if errors.Cause(err) == model.ErrProjectValidateFailed {
			logger.Info("Invalid project entry is specified: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to create project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	res := ProjectGetResponse{
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		TokenConfig: TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     project.TokenConfig.SigningAlgorithm,
		},
	}

	jwthttp.ResponseWrite(w, "ProjectCreateHandler", &res)
}

// ProjectDeleteHandler ...
//   require role: cluster-write
func ProjectDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r, role.ResCluster, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]

	if err := db.GetInst().ProjectDelete(projectName); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrDeleteBlockedProject {
			logger.Info("Failed to delete blocked project: %v", err)
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			logger.Error("Failed to delete project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("ProjectDeleteHandler method successfully finished")
}

// ProjectGetHandler ...
//   require role: project-read
func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Get Project
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	res := ProjectGetResponse{
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		TokenConfig: TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     project.TokenConfig.SigningAlgorithm,
		},
	}

	jwthttp.ResponseWrite(w, "ProjectGetHandler", &res)
}

// ProjectUpdateHandler ...
//   require role: project-write
func ProjectUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Parse Request
	var request ProjectPutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode project update request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Previous Project Info
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrProjectAlreadyExists {
			logger.Info("Project %s already exists", projectName)
			http.Error(w, "Project Already Exists", http.StatusBadRequest)
		} else {
			logger.Error("Failed to get project: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	project.TokenConfig.AccessTokenLifeSpan = request.TokenConfig.AccessTokenLifeSpan
	project.TokenConfig.RefreshTokenLifeSpan = request.TokenConfig.RefreshTokenLifeSpan
	project.TokenConfig.SigningAlgorithm = request.TokenConfig.SigningAlgorithm

	// Update DB
	if err := db.GetInst().ProjectUpdate(project); err != nil {
		logger.Error("Failed to update project: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("ProjectUpdateHandler method successfully finished")
}
