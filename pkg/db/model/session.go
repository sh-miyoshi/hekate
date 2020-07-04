package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/sh-miyoshi/hekate/pkg/errors"
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
	Add(projectName string, ent *Session) *errors.Error
	Delete(projectName string, sessionID string) *errors.Error
	DeleteAll(projectName string, userID string) *errors.Error
	DeleteAllInProject(projectName string) *errors.Error
	Get(projectName string, sessionID string) (*Session, *errors.Error)
	GetList(projectName string, userID string) ([]*Session, *errors.Error)
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
