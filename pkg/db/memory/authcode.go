package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
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
func (h *AuthCodeSessionHandler) Add(ent *model.AuthCodeSession) error {
	h.sessionList = append(h.sessionList, ent)
	return nil
}

// Update ...
func (h *AuthCodeSessionHandler) Update(ent *model.AuthCodeSession) error {
	for i, s := range h.sessionList {
		if s.SessionID == ent.SessionID {
			h.sessionList[i] = ent
			return nil
		}
	}
	return model.ErrNoSuchAuthCodeSession
}

// Delete ...
func (h *AuthCodeSessionHandler) Delete(sessionID string) error {
	newList := []*model.AuthCodeSession{}
	ok := false
	for _, s := range h.sessionList {
		if s.SessionID == sessionID {
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
func (h *AuthCodeSessionHandler) GetByCode(code string) (*model.AuthCodeSession, error) {
	for _, s := range h.sessionList {
		if s.Code == code {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchAuthCodeSession
}

// Get ...
func (h *AuthCodeSessionHandler) Get(id string) (*model.AuthCodeSession, error) {
	for _, s := range h.sessionList {
		if s.SessionID == id {
			return s, nil
		}
	}

	return nil, model.ErrNoSuchAuthCodeSession
}

// DeleteAllInClient ...
func (h *AuthCodeSessionHandler) DeleteAllInClient(clientID string) error {
	newList := []*model.AuthCodeSession{}
	for _, s := range h.sessionList {
		if s.ClientID != clientID {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

// DeleteAllInUser ...
func (h *AuthCodeSessionHandler) DeleteAllInUser(userID string) error {
	newList := []*model.AuthCodeSession{}
	for _, s := range h.sessionList {
		if s.UserID != userID {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}

// DeleteAllInProject ...
func (h *AuthCodeSessionHandler) DeleteAllInProject(projectName string) error {
	newList := []*model.AuthCodeSession{}
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}

	h.sessionList = newList
	return nil
}
