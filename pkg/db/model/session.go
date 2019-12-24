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
	New(userID string, sessionID string, expiresIn uint, fromIP string) error
	Revoke(sessionID string) error
	GetList(userID string) ([]string, error)
}
