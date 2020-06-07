package userapi

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
	"github.com/sh-miyoshi/hekate/pkg/pwpol"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

// AllUserGetHandler ...
//   require role: read-project
func AllUserGetHandler(w http.ResponseWriter, r *http.Request) {
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

	filter := &model.UserFilter{
		Name: queries.Get("name"),
	}

	users, err := db.GetInst().UserGetList(projectName, filter)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject || errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("Project %s is not found: %v", projectName, err)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Get all custom roles due to check all users
	customRoles, err := db.GetInst().CustomRoleGetList(projectName, nil)
	if err != nil {
		logger.Error("Failed to get custom role list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	res := []*UserGetResponse{}
	for _, user := range users {
		roles := []CustomRole{}
		for _, rid := range user.CustomRoles {
			for _, r := range customRoles {
				if rid == r.ID {
					roles = append(roles, CustomRole{
						r.ID,
						r.Name,
					})
					break
				}
			}
		}

		tmp := &UserGetResponse{
			ID:          user.ID,
			Name:        user.Name,
			CreatedAt:   user.CreatedAt.String(),
			SystemRoles: user.SystemRoles,
			CustomRoles: roles,
		}
		sessions, err := db.GetInst().SessionGetList(projectName, user.ID)
		if err != nil {
			logger.Error("Failed to get session list: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for _, s := range sessions {
			tmp.Sessions = append(tmp.Sessions, s.SessionID)
		}

		res = append(res, tmp)
	}

	jwthttp.ResponseWrite(w, "UserGetAllUserGetHandlerHandler", &res)
}

// UserCreateHandler ...
//   require role: write-project
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// validate password
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		logger.Error("Failed to get project: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := pwpol.CheckPassword(request.Name, request.Password, project.PasswordPolicy); err != nil {
		logger.Info("The password %s does not much the policy: %v", request.Password, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Create User Entry
	user := model.UserInfo{
		ID:           uuid.New().String(),
		ProjectName:  projectName,
		Name:         request.Name,
		CreatedAt:    time.Now(),
		PasswordHash: util.CreateHash(request.Password),
		SystemRoles:  request.SystemRoles,
		CustomRoles:  request.CustomRoles,
	}

	if err := db.GetInst().UserAdd(projectName, &user); err != nil {
		if errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("user validation failed: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrUserAlreadyExists {
			logger.Info("User %s is already exists", user.Name)
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			logger.Error("Failed to create user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	roles := []CustomRole{}
	for _, rid := range user.CustomRoles {
		r, err := db.GetInst().CustomRoleGet(projectName, rid)
		if err != nil {
			logger.Error("Failed to get user %s custom role %s info: %+v", user.ID, r.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		roles = append(roles, CustomRole{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	// Return Response
	res := UserGetResponse{
		ID:          user.ID,
		Name:        user.Name,
		CreatedAt:   user.CreatedAt.String(),
		SystemRoles: user.SystemRoles,
		CustomRoles: roles,
	}

	jwthttp.ResponseWrite(w, "UserGetAllUserGetHandlerHandler", &res)
}

// UserDeleteHandler ...
//   require role: write-project
func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete User
	if err := db.GetInst().UserDelete(projectName, userID); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchUser || errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("User %s is not found: %v", userID, err)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("UserDeleteHandler method successfully finished")
}

// UserGetHandler ...
//   require role: read-project, or <oneself>
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		claims, err := jwthttp.ValidateAPIRequest(r)
		// Check if the requester is the user
		if err != nil || claims.Subject != userID {
			logger.Info("Failed to authorize header: %v", err)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchUser || errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("User %s is not found: %v", userID, err)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("Invalid User ID format: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to get user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	roles := []CustomRole{}
	for _, rid := range user.CustomRoles {
		r, err := db.GetInst().CustomRoleGet(projectName, rid)
		if err != nil {
			logger.Error("Failed to get user %s custom role %s info: %+v", user.ID, r.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		roles = append(roles, CustomRole{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	res := UserGetResponse{
		ID:          user.ID,
		Name:        user.Name,
		CreatedAt:   user.CreatedAt.String(),
		SystemRoles: user.SystemRoles,
		CustomRoles: roles,
	}

	sessions, err := db.GetInst().SessionGetList(projectName, user.ID)
	if err != nil {
		logger.Error("Failed to get session list: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, s := range sessions {
		res.Sessions = append(res.Sessions, s.SessionID)
	}

	jwthttp.ResponseWrite(w, "UserGetHandler", &res)
}

// UserUpdateHandler ...
//   require role: write-project
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse Request
	var request UserPutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user update request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Previous User Info
	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchUser || errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("User %s is not found: %v", userID, err)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrUserValidateFailed {
			logger.Info("Invalid User ID format: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to update user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	// name, roles
	user.Name = request.Name
	user.SystemRoles = request.SystemRoles
	user.CustomRoles = request.CustomRoles

	// Update DB
	if err := db.GetInst().UserUpdate(projectName, user); err != nil {
		if errors.Cause(err) == model.ErrUserValidateFailed || errors.Cause(err) == model.ErrUserAlreadyExists {
			logger.Info("Invalid user request format: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to update user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("UserUpdateHandler method successfully finished")
}

// UserRoleAddHandler ...
//   require role: write-project
func UserRoleAddHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]
	roleID := vars["roleID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	roleType := model.RoleCustom
	if _, _, ok := role.GetInst().Parse(roleID); ok {
		roleType = model.RoleSystem
	}

	// Get Previous User Info
	if err := db.GetInst().UserAddRole(projectName, userID, roleType, roleID); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrRoleAlreadyAppended {
			logger.Info("Role %s is already appended", roleID)
			http.Error(w, "Role Already Appended", http.StatusConflict)
		} else if errors.Cause(err) == model.ErrUserValidateFailed {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				logger.Info("Invalid role was specified: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		} else {
			logger.Error("Failed to add role to user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserRoleAddHandler method successfully finished")
}

// UserRoleDeleteHandler ...
//   require role: write-project
func UserRoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]
	roleID := vars["roleID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get Previous User Info
	if err := db.GetInst().UserDeleteRole(projectName, userID, roleID); err != nil {
		if errors.Cause(err) == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrNoSuchRoleInUser {
			logger.Info("User %s do not have Role %s", userID, roleID)
			http.Error(w, "No Such Role in User", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrUserValidateFailed {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				logger.Info("Invalid ID was specified: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		} else {
			logger.Error("Failed to delete role from user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserRoleDeleteHandler method successfully finished")
}

// UserChangePasswordHandler ...
//   require role: <oneself>
func UserChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	claims, err := jwthttp.ValidateAPIRequest(r)
	if err != nil || claims.Subject != userID {
		logger.Info("Failed to authorize user: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UserChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info("Failed to decode user change password request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := db.GetInst().UserChangePassword(projectName, userID, req.Password); err != nil {
		if errors.Cause(err) == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if errors.Cause(err) == model.ErrUserValidateFailed {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				logger.Info("Invalid password was specified: %v", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}
		} else if errors.Cause(err) == pwpol.ErrPasswordPolicyFailed {
			logger.Info("Invalid password was specified: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			logger.Error("Failed to change user password: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserChangePasswordHandler method successfully finished")
}
