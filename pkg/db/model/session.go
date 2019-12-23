package model

import (
	"time"
)

// Session ...
type Session struct {
	UserID    string
	SessionID string
	CreatedAt time.Time
	ExpiresIn uint
	FromIP    string // Used to identify the user using this session
}

// SessionHandler ...
type SessionHandler interface {
	NewSession(userID string, sessionID string, expiresIn uint, fromIP string) error
	RevokeSession(sessionID string) error
	GetSessions(userID string) ([]string, error)
}
