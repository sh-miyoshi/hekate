package model

import (
	"time"

	"github.com/pkg/errors"
)

// AuthCodeSession ...
type AuthCodeSession struct {
	SessionID    string
	Code         string
	ExpiresIn    time.Time
	Scope        string
	ResponseType []string
	ClientID     string
	RedirectURI  string
	Nonce        string
	ProjectName  string
	MaxAge       uint
	ResponseMode string
	Prompt       string
	UserID       string
}

var (
	// ErrNoSuchAuthCodeSession ...
	ErrNoSuchAuthCodeSession = errors.New("No such session")
)

// AuthCodeSessionHandler ...
type AuthCodeSessionHandler interface {
	Add(ent *AuthCodeSession) error
	Update(ent *AuthCodeSession) error
	Delete(sessionID string) error
	GetByCode(code string) (*AuthCodeSession, error)
	Get(sessionID string) (*AuthCodeSession, error)
	DeleteAllInClient(clientID string) error
	DeleteAllInUser(userID string) error
	DeleteAllInProject(projectName string) error
}
