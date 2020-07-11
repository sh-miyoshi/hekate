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
func (h *LoginSessionHandler) Delete(projectName string, sessionID string) *errors.Error {
	newList := []*model.LoginSession{}
	ok := false
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == sessionID {
			ok = true
		} else {
			newList = append(newList, s)
		}
	}

	if !ok {
		return model.ErrNoSuchLoginSession
	}

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

// DeleteAllInClient ...
func (h *LoginSessionHandler) DeleteAllInClient(projectName string, clientID string) *errors.Error {
	newList := []*model.LoginSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		} else if s.ClientID != clientID {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

// DeleteAllInUser ...
func (h *LoginSessionHandler) DeleteAllInUser(projectName string, userID string) *errors.Error {
	newList := []*model.LoginSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		} else if s.UserID != userID {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

// DeleteAllInProject ...
func (h *LoginSessionHandler) DeleteAllInProject(projectName string) *errors.Error {
	newList := []*model.LoginSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}
