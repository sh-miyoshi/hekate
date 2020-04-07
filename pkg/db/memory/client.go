package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	clientList []*model.ClientInfo
}

// NewClientHandler ...
func NewClientHandler() *ClientInfoHandler {
	res := &ClientInfoHandler{}
	return res
}

// Add ...
func (h *ClientInfoHandler) Add(ent *model.ClientInfo) error {
	h.clientList = append(h.clientList, ent)
	return nil
}

// Delete ...
func (h *ClientInfoHandler) Delete(projectName, clientID string) error {
	newList := []*model.ClientInfo{}
	found := false
	for _, c := range h.clientList {
		if c.ProjectName == projectName && c.ID == clientID {
			found = true
		} else {
			newList = append(newList, c)
		}
	}

	if found {
		h.clientList = newList
		return nil
	}
	return model.ErrNoSuchClient
}

// GetList ...
func (h *ClientInfoHandler) GetList(projectName string) ([]*model.ClientInfo, error) {
	res := []*model.ClientInfo{}

	for _, client := range h.clientList {
		if client.ProjectName == projectName {
			res = append(res, client)
		}
	}

	return res, nil
}

// Get ...
func (h *ClientInfoHandler) Get(projectName, clientID string) (*model.ClientInfo, error) {
	for _, c := range h.clientList {
		if c.ProjectName == projectName && c.ID == clientID {
			return c, nil
		}
	}

	return nil, model.ErrNoSuchClient
}

// Update ...
func (h *ClientInfoHandler) Update(ent *model.ClientInfo) error {
	for i, c := range h.clientList {
		if c.ProjectName == ent.ProjectName && c.ID == ent.ID {
			h.clientList[i] = ent
			return nil
		}
	}
	return model.ErrNoSuchClient
}

// DeleteAll ...
func (h *ClientInfoHandler) DeleteAll(projectName string) error {
	newList := []*model.ClientInfo{}
	for _, c := range h.clientList {
		if c.ProjectName != projectName {
			newList = append(newList, c)
		}
	}
	h.clientList = newList
	return nil
}
