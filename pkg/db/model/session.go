package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

// Session ...
type Session struct {
	UserID      string
	ProjectName string
	SessionID   string
	CreatedAt   time.Time
	ExpiresIn   uint
	FromIP      string // Used to identify the user using this session
}

// SessionHandler ...
type SessionHandler interface {
	New(ent *Session) error
	Revoke(sessionID string) error
	RevokeAll(userID string) error
	Get(sessionID string) (*Session, error)
	GetList(userID string) ([]*Session, error)
}

var (
	// ErrSessionAlreadyExists ...
	ErrSessionAlreadyExists = errors.New("Session Already Exists")

	// ErrNoSuchSession ...
	ErrNoSuchSession = errors.New("No Such Session")

	// ErrSessionValidateFailed ...
	ErrSessionValidateFailed = errors.New("Session validation failed")
)

// Validate ...
func (s *Session) Validate() error {
	// Check Session ID
	if !ValidateSessionID(s.SessionID) {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid session ID format")
	}

	if !ValidateProjectName(s.ProjectName) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid project Name format")
	}

	// Check User ID
	if !ValidateUserID(s.UserID) {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid user ID format")
	}

	// Check From IP
	if ok := govalidator.IsIP(s.FromIP); !ok {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid from IP")
	}

	return nil
}
