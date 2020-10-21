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
	return errors.New("Internal Error", "No such client %s", clientID)
}

// GetList ...
func (h *ClientInfoHandler) GetList(projectName string, filter *model.ClientFilter) ([]*model.ClientInfo, *errors.Error) {
	res := []*model.ClientInfo{}

	for _, client := range h.clientList {
		if client.ProjectName == projectName {
			res = append(res, client)
		}
	}

	if filter != nil {
		res = matchFilterClientList(res, projectName, filter)
	}

	return res, nil
}

// Update ...
func (h *ClientInfoHandler) Update(projectName string, ent *model.ClientInfo) *errors.Error {
	for i, c := range h.clientList {
		if c.ProjectName == projectName && c.ID == ent.ID {
			h.clientList[i] = ent
			return nil
		}
	}
	return errors.New("Internal Error", "No such client %s", ent.ID)
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

// matchFilterClientList returns a list which matches the filter rules
func matchFilterClientList(data []*model.ClientInfo, projectName string, filter *model.ClientFilter) []*model.ClientInfo {
	if filter == nil {
		return data
	}
	res := []*model.ClientInfo{}

	for _, cli := range data {
		if projectName == cli.ProjectName {
			if filter.ID != "" && cli.ID != filter.ID {
				// missmatch id
				continue
			}
		}
		res = append(res, cli)
	}

	return res
}
