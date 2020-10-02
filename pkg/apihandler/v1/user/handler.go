package userapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/secret"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

// AllUserGetHandler ...
//   require role: read-project
func AllUserGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	queries := r.URL.Query()
	logger.Debug("Query: %v", queries)

	filter := &model.UserFilter{
		Name: queries.Get("name"),
	}

	users, err := db.GetInst().UserGetList(projectName, filter)
	if err != nil {
		if errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "User Get List Failed"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to get user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	// Get all custom roles due to check all users
	customRoles, err := db.GetInst().CustomRoleGetList(projectName, nil)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get custom role list"))
		errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			SystemRoles: user.SystemRoles,
			CustomRoles: roles,
			Locked:      user.LockState.Locked,
		}
		sessions, err := db.GetInst().SessionGetList(projectName, &model.SessionFilter{UserID: user.ID})
		if err != nil {
			errors.Print(errors.Append(err, "Failed to get session list"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	// Parse Request
	var request UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user create request: %v", err)
		errors.WriteHTTPError(w, "Bad Request", errors.New("Failed to decode request", ""), http.StatusBadRequest)
		return
	}

	// validate password
	project, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get project"))
		errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		return
	}

	if err := secret.CheckPassword(request.Name, request.Password, project.PasswordPolicy); err != nil {
		errors.PrintAsInfo(errors.Append(err, "The password does not much the policy"))
		errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
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
		if errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "user validation failed"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else if errors.Contains(err, model.ErrUserAlreadyExists) {
			errors.PrintAsInfo(errors.Append(err, "User %s is already exists", user.Name))
			errors.WriteHTTPError(w, "Conflict", err, http.StatusConflict)
		} else {
			errors.Print(errors.Append(err, "Failed to create user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	roles := []CustomRole{}
	for _, rid := range user.CustomRoles {
		r, err := db.GetInst().CustomRoleGet(projectName, rid)
		if err != nil {
			errors.Print(errors.Append(err, "Failed to get user %s custom role %s info", user.ID, r.ID))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		SystemRoles: user.SystemRoles,
		CustomRoles: roles,
		Locked:      user.LockState.Locked,
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	// Delete User
	if err := db.GetInst().UserDelete(projectName, userID); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) || errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "User %s is not found", userID))
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to delete user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
			errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
			errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
			return
		}
	}

	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) || errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "User %s is not found", userID))
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed toget user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	roles := []CustomRole{}
	for _, rid := range user.CustomRoles {
		r, err := db.GetInst().CustomRoleGet(projectName, rid)
		if err != nil {
			errors.Print(errors.Append(err, "Failed to get user %s custom role %s info", user.ID, r.ID))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		SystemRoles: user.SystemRoles,
		CustomRoles: roles,
		Locked:      user.LockState.Locked,
	}

	sessions, err := db.GetInst().SessionGetList(projectName, &model.SessionFilter{UserID: user.ID})
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get session list"))
		errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	// Parse Request
	var request UserPutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user update request: %v", err)
		errors.WriteHTTPError(w, "Bad Request", errors.New("Failed to decode request", ""), http.StatusBadRequest)
		return
	}

	// Get Previous User Info
	user, err := db.GetInst().UserGet(projectName, userID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			errors.PrintAsInfo(errors.Append(err, "User %s is not found", userID))
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Invalid user ID format"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to update user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		if errors.Contains(err, model.ErrUserValidateFailed) || errors.Contains(err, model.ErrUserAlreadyExists) {
			errors.PrintAsInfo(errors.Append(err, "Invalid user request format"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to update user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	roleType := model.RoleCustom
	if _, _, ok := role.GetInst().Parse(roleID); ok {
		roleType = model.RoleSystem
	}

	if err := db.GetInst().UserAddRole(projectName, userID, roleType, roleID); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			logger.Info("No such user: %s", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrRoleAlreadyAppended) {
			logger.Info("Role %s is already appended", roleID)
			errors.WriteHTTPError(w, "Conflict", err, http.StatusConflict)
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
			} else {
				// Includes role not found
				errors.PrintAsInfo(errors.Append(err, "Invalid role was specified"))
				errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
			}
		} else {
			errors.Print(errors.Append(err, "Failed to add role to user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	// Get Previous User Info
	if err := db.GetInst().UserDeleteRole(projectName, userID, roleID); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			logger.Info("No such user: %s", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrNoSuchRoleInUser) {
			logger.Info("User %s do not have Role %s", userID, roleID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
			} else {
				errors.PrintAsInfo(errors.Append(err, "Invalid ID was specified"))
				errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
			}
		} else {
			errors.Print(errors.Append(err, "Failed to delete role from user"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize user"))
		errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
		return
	}

	var req UserChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Info("Failed to decode user change password request: %v", err)
		errors.WriteHTTPError(w, "Bad Request", errors.New("Failed to decode request", ""), http.StatusBadRequest)
		return
	}

	if err := db.GetInst().UserChangePassword(projectName, userID, req.Password); err != nil {
		if errors.Contains(err, model.ErrNoSuchUser) {
			logger.Info("No such user: %s", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else if errors.Contains(err, model.ErrUserValidateFailed) {
			if !model.ValidateUserID(userID) {
				logger.Info("UserID %s is invalid id format", userID)
				errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
			} else {
				errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
				errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
			}
		} else if errors.Contains(err, secret.ErrPasswordPolicyFailed) {
			errors.PrintAsInfo(errors.Append(err, "Invalid password was specified"))
			errors.WriteHTTPError(w, "Bad Request", err, http.StatusBadRequest)
		} else {
			errors.Print(errors.Append(err, "Failed to change yser password"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserChangePasswordHandler method successfully finished")
}

// UserLogoutHandler ...
//   require role: <oneself> or write-project
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		claims, err := jwthttp.ValidateAPIRequest(r)
		if err != nil || claims.Subject != userID {
			errors.PrintAsInfo(errors.Append(err, "Failed to authorize user: don't have permission and not yourself"))
			errors.WriteHTTPError(w, "Forbidden", err, http.StatusForbidden)
			return
		}
	}

	if err := db.GetInst().UserLogout(projectName, userID); err != nil {
		if errors.Contains(err, model.ErrUserValidateFailed) {
			logger.Info("User ID %s is invalid", userID)
			errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
		} else {
			errors.Print(errors.Append(err, "Failed to logout"))
			errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserLogoutHandler method successfully finished")
}
