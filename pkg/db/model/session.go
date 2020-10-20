package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Session ...
type Session struct {
	UserID       string
	ProjectName  string
	SessionID    string
	CreatedAt    time.Time
	ExpiresIn    int64
	FromIP       string // Used to identify the user using this session
	LastAuthTime time.Time
}

// SessionFilter ...
type SessionFilter struct {
	SessionID string
	UserID    string
}

// SessionHandler ...
type SessionHandler interface {
	Add(projectName string, ent *Session) *errors.Error
	Delete(projectName string, filter *SessionFilter) *errors.Error
	DeleteAll(projectName string) *errors.Error
	GetList(projectName string, filter *SessionFilter) ([]*Session, *errors.Error)
}

var (
	// ErrSessionAlreadyExists ...
	ErrSessionAlreadyExists = errors.New("Session already exists", "Session already exists")

	// ErrNoSuchSession ...
	ErrNoSuchSession = errors.New("No such session", "No such session")

	// ErrSessionValidateFailed ...
	ErrSessionValidateFailed = errors.New("Session validation failed", "Session validation failed")
)

// Validate ...
func (s *Session) Validate() *errors.Error {
	// Check Session ID
	if !ValidateSessionID(s.SessionID) {
		return errors.Append(ErrSessionValidateFailed, "Invalid session ID format")
	}

	if !ValidateProjectName(s.ProjectName) {
		return errors.Append(ErrUserValidateFailed, "Invalid project Name format")
	}

	// Check User ID
	if !ValidateUserID(s.UserID) {
		return errors.Append(ErrSessionValidateFailed, "Invalid user ID format")
	}

	// Check From IP
	if ok := govalidator.IsIP(s.FromIP); !ok {
		return errors.Append(ErrSessionValidateFailed, "Invalid from IP")
	}

	return nil
}
