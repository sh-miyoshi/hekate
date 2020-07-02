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
	Prompt       []string
	UserID       string
	LoginDate    time.Time
}

var (
	// ErrNoSuchAuthCodeSession ...
	ErrNoSuchAuthCodeSession = errors.New("No such session")
)

// AuthCodeSessionHandler ...
type AuthCodeSessionHandler interface {
	Add(projectName string, ent *AuthCodeSession) error
	Update(projectName string, ent *AuthCodeSession) error
	Delete(projectName string, sessionID string) error
	GetByCode(projectName string, code string) (*AuthCodeSession, error)
	Get(projectName string, sessionID string) (*AuthCodeSession, error)
	DeleteAllInClient(projectName string, clientID string) error
	DeleteAllInUser(projectName string, userID string) error
	DeleteAllInProject(projectName string) error
}
