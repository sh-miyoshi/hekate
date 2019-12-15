package userapi

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/jwt-server/pkg/http"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/role"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"time"
)

// AllUserGetHandler ...
//   require role: project-read
func AllUserGetHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]

	users, err := db.GetInst().User.GetList(projectName)
	if err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	jwthttp.ResponseWrite(w, "UserGetAllUserGetHandlerHandler", &users)
}

// UserCreateHandler ...
//   require role: project-write
func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Parse Request
	var request UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user create request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Validate Request)

	// Create User Entry
	user := model.UserInfo{
		ID:           uuid.New().String(),
		ProjectName:  projectName,
		Name:         request.Name,
		CreatedAt:    time.Now(),
		PasswordHash: util.CreateHash(request.Password),
		Roles:        request.Roles,
	}

	if err := db.GetInst().User.Add(&user); err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if err == model.ErrUserAlreadyExists {
			logger.Info("User %s is already exists", user.Name)
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			logger.Error("Failed to create user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return Response
	res := UserGetResponse{
		ID:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		Roles:        user.Roles,
	}

	jwthttp.ResponseWrite(w, "UserGetAllUserGetHandlerHandler", &res)

	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}

// UserDeleteHandler ...
//   require role: project-write
func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	if err := db.GetInst().User.Delete(projectName, userID); err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if err == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("PUserDeleteHandler method successfully finished")
}

// UserGetHandler ...
//   require role: user-read
func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResUser, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	user, err := db.GetInst().User.Get(projectName, userID)
	if err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to get user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	res := UserGetResponse{
		ID:           user.ID,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		Roles:        user.Roles,
	}

	for _, s := range user.Sessions {
		res.Sessions = append(res.Sessions, s.SessionID)
	}

	jwthttp.ResponseWrite(w, "UserGetHandler", &res)
}

// UserUpdateHandler ...
//   require role: user-write
func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResUser, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]

	// Parse Request
	var request UserPutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Info("Failed to decode user update request: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Previous User Info
	user, err := db.GetInst().User.Get(projectName, userID)
	if err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if err == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else {
			logger.Error("Failed to update user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update Parameters
	// name, password, roles
	user.Name = request.Name
	user.PasswordHash = util.CreateHash(request.Password)
	user.Roles = request.Roles

	// Update DB
	if err := db.GetInst().User.Update(user); err != nil {
		logger.Error("Failed to update user: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("UserUpdateHandler method successfully finished")
}

// UserRoleAddHandler ...
//   require role: user-write
func UserRoleAddHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResUser, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]
	roleID := vars["roleID"]

	// Get Previous User Info
	if err := db.GetInst().User.AppendRole(projectName, userID, roleID); err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if err == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if err == model.ErrRoleAlreadyAppended {
			logger.Info("Role %s is already appended", roleID)
			http.Error(w, "Role Already Appended", http.StatusConflict)
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
//   require role: user-write
func UserRoleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize API Request
	if err := jwthttp.AuthHeader(r.Header, role.ResUser, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	projectName := vars["projectName"]
	userID := vars["userID"]
	roleID := vars["roleID"]

	// Get Previous User Info
	if err := db.GetInst().User.RemoveRole(projectName, userID, roleID); err != nil {
		if err == model.ErrNoSuchProject {
			logger.Info("No such project: %s", projectName)
			http.Error(w, "Project Not Found", http.StatusNotFound)
		} else if err == model.ErrNoSuchUser {
			logger.Info("No such user: %s", userID)
			http.Error(w, "User Not Found", http.StatusNotFound)
		} else if err == model.ErrNoSuchRoleInUser {
			logger.Info("User %s do not have Role %s", userID, roleID)
			http.Error(w, "No Such Role in User", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete role from user: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("UserRoleDeleteHandler method successfully finished")

	logger.Info("Not implemented yet")
	http.Error(w, "Not Implemented yet", http.StatusInternalServerError)
}
