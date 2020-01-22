package memory

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"sync"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	// sessionList[sessionID] = Session
	sessionList map[string]*model.Session
	mu          sync.Mutex
}

// NewSessionHandler ...
func NewSessionHandler() (*SessionHandler, error) {
	res := &SessionHandler{
		sessionList: make(map[string]*model.Session),
	}
	return res, nil
}

// New ...
func (h *SessionHandler) New(session *model.Session) error {
	h.sessionList[session.SessionID] = session
	return nil
}

// Revoke ...
func (h *SessionHandler) Revoke(sessionID string) error {
	if _, exists := h.sessionList[sessionID]; exists {
		delete(h.sessionList, sessionID)
		return nil
	}

	return model.ErrNoSuchSession
}

// Get ...
func (h *SessionHandler) Get(sessionID string) (*model.Session, error) {
	if _, exists := h.sessionList[sessionID]; !exists {
		return nil, model.ErrNoSuchSession
	}
	return h.sessionList[sessionID], nil
}

// GetList ...
func (h *SessionHandler) GetList(userID string) ([]string, error) {
	res := []string{}

	for id, s := range h.sessionList {
		if s.UserID == userID {
			res = append(res, id)
		}
	}

	return res, nil
}

// BeginTx ...
func (h *SessionHandler) BeginTx() error {
	h.mu.Lock()
	return nil
}

// CommitTx ...
func (h *SessionHandler) CommitTx() error {
	h.mu.Unlock()
	return nil
}

// AbortTx ...
func (h *SessionHandler) AbortTx() error {
	h.mu.Unlock()
	return nil
}
