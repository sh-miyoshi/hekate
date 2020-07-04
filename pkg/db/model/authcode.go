package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
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
	Add(projectName string, ent *AuthCodeSession) *errors.Error
	Update(projectName string, ent *AuthCodeSession) *errors.Error
	Delete(projectName string, sessionID string) *errors.Error
	GetByCode(projectName string, code string) (*AuthCodeSession, *errors.Error)
	Get(projectName string, sessionID string) (*AuthCodeSession, *errors.Error)
	DeleteAllInClient(projectName string, clientID string) *errors.Error
	DeleteAllInUser(projectName string, userID string) *errors.Error
	DeleteAllInProject(projectName string) *errors.Error
}
