package model

import (
	"errors"
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
	Get(sessionID string) (*Session, error)
	GetList(userID string) ([]string, error)
}

var (
	// ErrSessionAlreadyExists ...
	ErrSessionAlreadyExists = errors.New("Session Already Exists")

	// ErrNoSuchSession ...
	ErrNoSuchSession = errors.New("No Such Session")
)
