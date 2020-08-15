package customroleapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
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
		if errors.Contains(err, model.ErrCustomRoleValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Failed to get role list"))
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to get role list"))
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

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "ROLE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
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
		if errors.Contains(err, model.ErrCustomRoleAlreadyExists) {
			errors.PrintAsInfo(errors.Append(err, "Custom Role %s is already exists", role.Name))
			http.Error(w, "Custom Role already exists", http.StatusConflict)
		} else if errors.Contains(err, model.ErrCustomRoleValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Custom role validation failed"))
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to create role"))
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

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "ROLE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := db.GetInst().CustomRoleDelete(projectName, roleID); err != nil {
		if errors.Contains(err, model.ErrNoSuchCustomRole) || errors.Contains(err, model.ErrCustomRoleValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Custom role %s is not found", roleID))
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to delete custom role"))
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	role, err := db.GetInst().CustomRoleGet(projectName, roleID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchCustomRole) || errors.Contains(err, model.ErrCustomRoleValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Custom role %s is not found", roleID))
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to get role"))
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

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "ROLE", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
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
		if errors.Contains(err, model.ErrNoSuchCustomRole) || errors.Contains(err, model.ErrCustomRoleValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Custom role %s is not found", request.Name))
			http.Error(w, "Custom Role Not Found", http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to update role"))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	role.Name = request.Name

	// Update DB
	if err := db.GetInst().CustomRoleUpdate(projectName, role); err != nil {
		if errors.Contains(err, model.ErrCustomRoleValidateFailed) || errors.Contains(err, model.ErrCustomRoleAlreadyExists) {
			errors.PrintAsInfo(errors.Append(err, "Failed to validate request"))
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to update role"))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("CustomRoleUpdateHandler method successfully finished")
}
