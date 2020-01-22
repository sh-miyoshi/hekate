package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
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
	New(ent *Session) error
	Revoke(sessionID string) error
	Get(sessionID string) (*Session, error)
	GetList(userID string) ([]string, error)

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
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
	if ok := govalidator.IsUUID(s.SessionID); !ok {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid session ID format")
	}

	// Check User ID
	if ok := govalidator.IsUUID(s.UserID); !ok {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid user ID format")
	}

	// Check From IP
	if ok := govalidator.IsIP(s.FromIP); !ok {
		return errors.Wrap(ErrSessionValidateFailed, "Invalid from IP")
	}

	return nil
}
