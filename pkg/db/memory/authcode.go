package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// AuthCodeSessionHandler implement db.AuthCodeSessionHandler
type AuthCodeSessionHandler struct {
	sessionList []*model.AuthCodeSession
}

// NewAuthCodeSessionHandler ...
func NewAuthCodeSessionHandler() *AuthCodeSessionHandler {
	return &AuthCodeSessionHandler{}
}

// Add ...
func (h *AuthCodeSessionHandler) Add(projectName string, ent *model.AuthCodeSession) *errors.Error {
	h.sessionList = append(h.sessionList, ent)
	return nil
}

// Update ...
func (h *AuthCodeSessionHandler) Update(projectName string, ent *model.AuthCodeSession) *errors.Error {
	for i, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == ent.SessionID {
			h.sessionList[i] = ent
			return nil
		}
	}
	return model.ErrNoSuchAuthCodeSession
}

// Delete ...
func (h *AuthCodeSessionHandler) Delete(projectName string, sessionID string) *errors.Error {
	newList := []*model.AuthCodeSession{}
	ok := false
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == sessionID {
			ok = true
		} else {
			newList = append(newList, s)
		}
	}

	if !ok {
		return model.ErrNoSuchAuthCodeSession
	}

	h.sessionList = newList
	return nil
}

// GetByCode ...
func (h *AuthCodeSessionHandler) GetByCode(projectName string, code string) (*model.AuthCodeSession, *errors.Error) {
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.Code == code {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchAuthCodeSession
}

// Get ...
func (h *AuthCodeSessionHandler) Get(projectName string, id string) (*model.AuthCodeSession, *errors.Error) {
	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.SessionID == id {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchAuthCodeSession
}

// DeleteAllInClient ...
func (h *AuthCodeSessionHandler) DeleteAllInClient(projectName string, clientID string) *errors.Error {
	newList := []*model.AuthCodeSession{}
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
func (h *AuthCodeSessionHandler) DeleteAllInUser(projectName string, userID string) *errors.Error {
	newList := []*model.AuthCodeSession{}
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
func (h *AuthCodeSessionHandler) DeleteAllInProject(projectName string) *errors.Error {
	newList := []*model.AuthCodeSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}
