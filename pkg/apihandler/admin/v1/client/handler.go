package clientapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	jwthttp "github.com/sh-miyoshi/hekate/pkg/http"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
)

// AllClientGetHandler ...
//   require role: read-project
func AllClientGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	clients, err := db.GetInst().ClientGetList(projectName, nil)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get client"))
		errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		return
	}

	res := []*ClientGetResponse{}
	for _, client := range clients {
		res = append(res, &ClientGetResponse{
			ID:                  client.ID,
			Secret:              client.Secret,
			AccessType:          client.AccessType,
			CreatedAt:           client.CreatedAt.Format(time.RFC3339),
			AllowedCallbackURLs: client.AllowedCallbackURLs,
		})
	}

	jwthttp.ResponseWrite(w, "AllClientGetHandlerHandler", res)
}

// ClientCreateHandler ...
//   require role: write-project
func ClientCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "CLIENT", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err = jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	// Parse Request
	var request ClientCreateRequest
	if e := json.NewDecoder(r.Body).Decode(&request); e != nil {
		err = errors.Append(errors.ErrInvalidRequest, "Failed to decode client create request: %v", e)
		errors.PrintAsInfo(err)
		errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		return
	}

	// Create Client Entry
	client := model.ClientInfo{
		ID:                  request.ID,
		ProjectName:         projectName,
		Secret:              request.Secret,
		AccessType:          request.AccessType,
		CreatedAt:           time.Now(),
		AllowedCallbackURLs: request.AllowedCallbackURLs,
	}

	if err = db.GetInst().ClientAdd(projectName, &client); err != nil {
		if errors.Contains(err, model.ErrClientAlreadyExists) {
			errors.PrintAsInfo(errors.Append(err, "Client %s is already exists", client.ID))
			errors.WriteToHTTP(w, err, http.StatusConflict, "")
		} else if errors.Contains(err, model.ErrClientValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Bad Request"))
			errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		} else {
			errors.Print(errors.Append(err, "Failed to create client"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	// Return Response
	res := ClientGetResponse{
		ID:                  client.ID,
		Secret:              client.Secret,
		AccessType:          client.AccessType,
		CreatedAt:           client.CreatedAt.Format(time.RFC3339),
		AllowedCallbackURLs: client.AllowedCallbackURLs,
	}

	jwthttp.ResponseWrite(w, "ClientCreateHandler", &res)
}

// ClientDeleteHandler ...
//   require role: write-project
func ClientDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	clientID := vars["clientID"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "CLIENT", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err = jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	if err = db.GetInst().ClientDelete(projectName, clientID); err != nil {
		if errors.Contains(err, model.ErrNoSuchClient) || errors.Contains(err, model.ErrClientValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "No such client: %s", clientID))
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else {
			errors.Print(errors.Append(err, "Failed to delete client"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	// Return 204 (No content) for success
	w.WriteHeader(http.StatusNoContent)
	logger.Info("ClientDeleteHandler method successfully finished")
}

// ClientGetHandler ...
//   require role: read-project
func ClientGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	clientID := vars["clientID"]

	// Authorize API Request
	if err := jwthttp.Authorize(r, projectName, role.ResProject, role.TypeRead); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	client, err := db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchClient) || errors.Contains(err, model.ErrClientValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "No such client: %s", clientID))
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else {
			errors.Print(errors.Append(err, "Failed to get client"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	res := ClientGetResponse{
		ID:                  client.ID,
		Secret:              client.Secret,
		AccessType:          client.AccessType,
		CreatedAt:           client.CreatedAt.Format(time.RFC3339),
		AllowedCallbackURLs: client.AllowedCallbackURLs,
	}

	jwthttp.ResponseWrite(w, "ClientGetHandler", &res)
}

// ClientUpdateHandler ...
//   require role: write-project
func ClientUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["projectName"]
	clientID := vars["clientID"]

	var err *errors.Error
	defer func() {
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		if err = audit.GetInst().Save(projectName, time.Now(), "CLIENT", r.Method, r.URL.String(), msg); err != nil {
			errors.Print(errors.Append(err, "Failed to save audit event"))
		}
	}()

	// Authorize API Request
	if err = jwthttp.Authorize(r, projectName, role.ResProject, role.TypeWrite); err != nil {
		errors.PrintAsInfo(errors.Append(err, "Failed to authorize header"))
		errors.WriteToHTTP(w, errors.ErrUnpermitted, 0, "")
		return
	}

	// Parse Request
	var request ClientPutRequest
	if e := json.NewDecoder(r.Body).Decode(&request); e != nil {
		err = errors.Append(errors.ErrInvalidRequest, "Failed to decode client update request: %v", e)
		errors.PrintAsInfo(err)
		errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		return
	}

	// Get Previous Client Info
	var client *model.ClientInfo
	client, err = db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchClient) || errors.Contains(err, model.ErrClientValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "No such client: %s", clientID))
			errors.WriteToHTTP(w, err, http.StatusNotFound, "")
		} else {
			errors.Print(errors.Append(err, "Failed to update client"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	// Update Parameters
	client.Secret = request.Secret
	client.AccessType = request.AccessType
	client.AllowedCallbackURLs = request.AllowedCallbackURLs

	// Update DB
	if err = db.GetInst().ClientUpdate(projectName, client); err != nil {
		if errors.Contains(err, model.ErrClientValidateFailed) {
			errors.PrintAsInfo(errors.Append(err, "Bad Request"))
			errors.WriteToHTTP(w, err, http.StatusBadRequest, "")
		} else {
			errors.Print(errors.Append(err, "Failed to update client"))
			errors.WriteToHTTP(w, err, http.StatusInternalServerError, "")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("ClientUpdateHandler method successfully finished")
}
