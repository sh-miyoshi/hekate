package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
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
func (h *ClientInfoHandler) Add(projectName string, ent *model.ClientInfo) *errors.Error {
	h.clientList = append(h.clientList, ent)
	return nil
}

// Delete ...
func (h *ClientInfoHandler) Delete(projectName, clientID string) *errors.Error {
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
func (h *ClientInfoHandler) GetList(projectName string) ([]*model.ClientInfo, *errors.Error) {
	res := []*model.ClientInfo{}

	for _, client := range h.clientList {
		if client.ProjectName == projectName {
			res = append(res, client)
		}
	}

	return res, nil
}

// Get ...
func (h *ClientInfoHandler) Get(projectName, clientID string) (*model.ClientInfo, *errors.Error) {
	for _, c := range h.clientList {
		if c.ProjectName == projectName && c.ID == clientID {
			return c, nil
		}
	}

	return nil, model.ErrNoSuchClient
}

// Update ...
func (h *ClientInfoHandler) Update(projectName string, ent *model.ClientInfo) *errors.Error {
	for i, c := range h.clientList {
		if c.ProjectName == projectName && c.ID == ent.ID {
			h.clientList[i] = ent
			return nil
		}
	}
	return model.ErrNoSuchClient
}

// DeleteAll ...
func (h *ClientInfoHandler) DeleteAll(projectName string) *errors.Error {
	newList := []*model.ClientInfo{}
	for _, c := range h.clientList {
		if c.ProjectName != projectName {
			newList = append(newList, c)
		}
	}
	h.clientList = newList
	return nil
}
