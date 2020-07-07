package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	// sessionList[sessionID] = Session
	sessionList map[string]*model.Session
}

// NewSessionHandler ...
func NewSessionHandler() *SessionHandler {
	res := &SessionHandler{
		sessionList: make(map[string]*model.Session),
	}
	return res
}

// Add ...
func (h *SessionHandler) Add(projectName string, session *model.Session) *errors.Error {
	h.sessionList[session.SessionID] = session
	return nil
}

// Delete ...
func (h *SessionHandler) Delete(projectName string, sessionID string) *errors.Error {
	if res, exists := h.sessionList[sessionID]; exists {
		if res.ProjectName == projectName {
			delete(h.sessionList, sessionID)
			return nil
		}
	}

	return model.ErrNoSuchSession
}

// DeleteAll ...
func (h *SessionHandler) DeleteAll(projectName string, userID string) *errors.Error {
	newList := make(map[string]*model.Session)
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList[s.SessionID] = s
		} else if s.UserID != userID {
			newList[s.SessionID] = s
		}
	}
	h.sessionList = newList
	return nil
}

// DeleteAllInProject ...
func (h *SessionHandler) DeleteAllInProject(projectName string) *errors.Error {
	newList := make(map[string]*model.Session)
	for _, s := range h.sessionList {
		if s.ProjectName != projectName {
			newList[s.SessionID] = s
		}
	}
	h.sessionList = newList
	return nil
}

// Get ...
func (h *SessionHandler) Get(projectName string, sessionID string) (*model.Session, *errors.Error) {
	res, exists := h.sessionList[sessionID]
	if !exists || res.ProjectName != projectName {
		return nil, model.ErrNoSuchSession
	}
	return res, nil
}

// GetList ...
func (h *SessionHandler) GetList(projectName string, userID string) ([]*model.Session, *errors.Error) {
	res := []*model.Session{}

	for _, s := range h.sessionList {
		if s.ProjectName == projectName && s.UserID == userID {
			res = append(res, s)
		}
	}

	return res, nil
}
