package sessionapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// SessionDeleteHandler ...
//   require role: write-project
func SessionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	sessionID := vars["sessionID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := db.GetInst().SessionDelete(projectName, sessionID); err != nil {
		e := errors.Cause(err)
		if e == model.ErrNoSuchProject || e == model.ErrNoSuchSession || e == model.ErrSessionValidateFailed {
			logger.Error("Failed to delete session: %v", err)
			http.Error(w, "No such session", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete session info: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("SessionDeleteHandler method successfully finished")
}

// SessionGetHandler ...
//   require role: read-project
func SessionGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	sessionID := vars["sessionID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		logger.Info("Failed to authorize header: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	s, err := db.GetInst().SessionGet(projectName, sessionID)
	if err != nil {
		e := errors.Cause(err)
		if e == model.ErrNoSuchProject || e == model.ErrNoSuchSession || e == model.ErrSessionValidateFailed {
			logger.Error("Failed to get session: %v", err)
			http.Error(w, "No such session", http.StatusNotFound)
		} else {
			logger.Error("Failed to get session info: %+v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	res := SessionGetResponse{
		ID:        s.SessionID,
		CreatedAt: s.CreatedAt.String(),
		ExpiresIn: s.ExpiresIn,
		FromIP:    s.FromIP,
	}

	jwthttp.ResponseWrite(w, "SessionGetHandler", &res)
}
