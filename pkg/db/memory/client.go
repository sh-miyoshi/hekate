package memory

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"sync"
)

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	// clientList[clientID] = ClientInfo
	clientList     map[string]*model.ClientInfo
	projectHandler *ProjectInfoHandler
	mu             sync.Mutex
}

// NewClientHandler ...
func NewClientHandler(projectHandler *ProjectInfoHandler) (*ClientInfoHandler, error) {
	res := &ClientInfoHandler{
		clientList:     make(map[string]*model.ClientInfo),
		projectHandler: projectHandler,
	}
	return res, nil
}

// Add ...
func (h *ClientInfoHandler) Add(ent *model.ClientInfo) error {
	h.clientList[ent.ID] = ent
	return nil
}

// Delete ...
func (h *ClientInfoHandler) Delete(clientID string) error {
	if _, exists := h.clientList[clientID]; exists {
		delete(h.clientList, clientID)
		return nil
	}
	return model.ErrNoSuchClient
}

// GetList ...
func (h *ClientInfoHandler) GetList(projectName string) ([]string, error) {
	res := []string{}

	if _, err := h.projectHandler.Get(projectName); err != nil {
		// project is created in Add method, so maybe empty project
		return res, nil
	}

	for _, client := range h.clientList {
		if client.ProjectName == projectName {
			res = append(res, client.ID)
		}
	}

	return res, nil
}

// Get ...
func (h *ClientInfoHandler) Get(clientID string) (*model.ClientInfo, error) {
	res, exists := h.clientList[clientID]
	if !exists {
		return nil, model.ErrNoSuchClient
	}

	return res, nil
}

// Update ...
func (h *ClientInfoHandler) Update(ent *model.ClientInfo) error {
	if _, exists := h.clientList[ent.ID]; !exists {
		return model.ErrNoSuchClient
	}

	h.clientList[ent.ID] = ent

	return nil
}

// DeleteAll ...
func (h *ClientInfoHandler) DeleteAll(projectName string) error {
	for _, client := range h.clientList {
		if client.ProjectName == projectName {
			delete(h.clientList, client.ID)
		}
	}
	return nil
}

// BeginTx ...
func (h *ClientInfoHandler) BeginTx() error {
	h.mu.Lock()
	return nil
}

// CommitTx ...
func (h *ClientInfoHandler) CommitTx() error {
	h.mu.Unlock()
	return nil
}

// AbortTx ...
func (h *ClientInfoHandler) AbortTx() error {
	h.mu.Unlock()
	return nil
}
