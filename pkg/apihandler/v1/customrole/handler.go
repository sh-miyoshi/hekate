package customroleapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// AllRoleGetHandler ...
//   require role: read-project
func AllRoleGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	filter := &model.CustomRoleFilter{
		Name: queries.Get("name"),
	}

	roles, err := db.GetInst().CustomRoleGetList(projectName, filter)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject || errors.Cause(err) == model.ErrCustomRoleValidateFailed {
			logger.Info("Failed to get role list: %v", err)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get role list: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	res := []*CustomRoleGetResponse{}
	for _, role := range roles {
		res = append(res, &CustomRoleGetResponse{
			ID:          role.ID,
			Name:        role.Name,
			ProjectName: role.ProjectName,
			CreatedAt:   role.CreatedAt.String(),
		})
	}

	jwthttp.ResponseWrite(w, "AllRoleGetHandler", res)

}

// RoleCreateHandler ...
//   require role: write-project
func RoleCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request CustomRoleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode role create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create role entry
	role := model.CustomRole{
		ID:          uuid.New().String(),
		Name:        request.Name,
		ProjectName: projectName,
		CreatedAt:   time.Now(),
	}

	if err := db.GetInst().CustomRoleAdd(projectName, &role); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrCustomRoleAlreadyExists {
			logger.Info("Custom Role %s is already exists", role.Name)
			http.Error(w, "Custom Role already exists", http.StatusConflict)
		} else if errors.Cause(err) == model.ErrCustomRoleValidateFailed {
			logger.Info("Custom role validation failed: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to create role: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	res := CustomRoleGetResponse{
		ID:          role.ID,
		Name:        role.Name,
		ProjectName: role.ProjectName,
		CreatedAt:   role.CreatedAt.String(),
	}

	jwthttp.ResponseWrite(w, "RoleCreateHandler", &res)

}

// RoleDeleteHandler ...
//   require role: write-project
func RoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	roleID := vars["roleID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := db.GetInst().CustomRoleDelete(projectName, roleID); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchCustomRole || errors.Cause(err) == model.ErrCustomRoleValidateFailed {
			logger.Info("Custom role %s is not found: %v", roleID, err)
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete custom role: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("RoleDeleteHandler method successfully finished")
}

// RoleGetHandler ...
//   require role: read-project
func RoleGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	roleID := vars["roleID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	role, err := db.GetInst().CustomRoleGet(projectName, roleID)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchCustomRole || errors.Cause(err) == model.ErrCustomRoleValidateFailed {
			logger.Info("Custom role %s is not found: %v", roleID, err)
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get role: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	res := CustomRoleGetResponse{
		ID:          role.ID,
		Name:        role.Name,
		ProjectName: role.ProjectName,
		CreatedAt:   role.CreatedAt.String(),
	}

	jwthttp.ResponseWrite(w, "CustomRoleGetHandler", &res)
}

// RoleUpdateHandler ...
//   require role: write-project
func RoleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	roleID := vars["roleID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request CustomRolePutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode role update request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Previous CustomRole Info
	role, err := db.GetInst().CustomRoleGet(projectName, roleID)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchCustomRole || errors.Cause(err) == model.ErrCustomRoleValidateFailed {
			logger.Info("Custom role %s is not found: %v", roleID, err)
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to update role: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	role.Name = request.Name

	// Update DB
	if err := db.GetInst().CustomRoleUpdate(projectName, role); err != nil {
		if errors.Cause(err) == model.ErrCustomRoleValidateFailed || errors.Cause(err) == model.ErrCustomRoleAlreadyExists {
			logger.Error("Request validation failed: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to update role: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("CustomRoleUpdateHandler method successfully finished")
}
