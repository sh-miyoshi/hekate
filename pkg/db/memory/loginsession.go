package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// LoginSessionHandler implement db.LoginSessionHandler
type LoginSessionHandler struct {
	sessionList []*model.LoginSession
}

// NewLoginSessionHandler ...
func NewLoginSessionHandler() *LoginSessionHandler {
	return &LoginSessionHandler{}
}

// Add ...
func (h *LoginSessionHandler) Add(projectName string, ent *model.LoginSession) *errors.Error {
	h.sessionList = append(h.sessionList, ent)
	return nil
}

// Update ...
func (h *LoginSessionHandler) Update(projectName string, ent *model.LoginSession) *errors.Error {
	for i, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == ent.SessionID {
			h.sessionList[i] = ent
			return nil
		}
	}
	return model.ErrNoSuchLoginSession
}

// Delete ...
func (h *LoginSessionHandler) Delete(projectName string, filter *model.LoginSessionFilter) *errors.Error {
	newList := filterLoginSessionList(h.sessionList, projectName, filter)

	h.sessionList = newList
	return nil
}

// GetByCode ...
func (h *LoginSessionHandler) GetByCode(projectName string, code string) (*model.LoginSession, *errors.Error) {
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.Code == code {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchLoginSession
}

// Get ...
func (h *LoginSessionHandler) Get(projectName string, id string) (*model.LoginSession, *errors.Error) {
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == id {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchLoginSession
}

// DeleteAll ...
func (h *LoginSessionHandler) DeleteAll(projectName string) *errors.Error {
	newList := []*model.LoginSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

func filterLoginSessionList(data []*model.LoginSession, projectName string, filter *model.LoginSessionFilter) []*model.LoginSession {
	if filter == nil {
		return data
	}
	res := []*model.LoginSession{}

	for _, s := range data {
		if projectName == s.ProjectName {
			if filter.SessionID != "" && s.SessionID != filter.SessionID {
				// missmatch session id
				continue
			}
			if filter.UserID != "" && s.UserID != filter.UserID {
				// missmatch user id
				continue
			}
			if filter.ClientID != "" && s.ClientID != filter.ClientID {
				// missmatch session id
				continue
			}
		}
		res = append(res, s)
	}

	return res
}
