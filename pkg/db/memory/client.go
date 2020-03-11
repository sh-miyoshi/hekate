package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	// clientList[clientID] = ClientInfo
	clientList map[string]*model.ClientInfo
}

// NewClientHandler ...
func NewClientHandler() *ClientInfoHandler {
	res := &ClientInfoHandler{
		clientList: make(map[string]*model.ClientInfo),
	}
	return res
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
