package memory

import (
	"time"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	// sessionList[sessionID] = Session
	sessionList       map[string]*model.Session
}

// NewSessionHandler ...
func NewSessionHandler() (*SessionHandler, error) {
	res := &SessionHandler{
		sessionList:       make(map[string]*model.Session),
	}
	return res, nil
}

// NewSession ...
func (h *SessionHandler) NewSession(userID string, sessionID string, expiresIn uint, fromIP string) error {
	if _, exists := h.sessionList[sessionID]; exists {
		return errors.New("Session already exists")
	}

	h.sessionList[sessionID] = &model.Session{
		UserID: userID, 
		SessionID: sessionID,
		CreatedAt: time.Now(),
		ExpiresIn: expiresIn,
		FromIP: fromIP,
	}

	return nil
}

// RevokeSession ...
func (h *SessionHandler) RevokeSession(sessionID string) error {
	if _, exists := h.sessionList[sessionID]; exists {
		delete(h.sessionList, sessionID)
		return nil
	}

	return nil
}

// GetSessions ...
func (h *SessionHandler) GetSessions(userID string) ([]string, error){
	res := []string{}

	for id, s := range h.sessionList {
		if s.UserID == userID {
			res = append(res, id)
		}
	}

	return res, nil
}