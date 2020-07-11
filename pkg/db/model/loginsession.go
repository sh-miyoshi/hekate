package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// LoginSession ...
type LoginSession struct {
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
	// ErrNoSuchLoginSession ...
	ErrNoSuchLoginSession = errors.New("No such session", "No such session")
)

// LoginSessionHandler ...
type LoginSessionHandler interface {
	Add(projectName string, ent *LoginSession) *errors.Error
	Update(projectName string, ent *LoginSession) *errors.Error
	Delete(projectName string, sessionID string) *errors.Error
	GetByCode(projectName string, code string) (*LoginSession, *errors.Error)
	Get(projectName string, sessionID string) (*LoginSession, *errors.Error)
	DeleteAllInClient(projectName string, clientID string) *errors.Error
	DeleteAllInUser(projectName string, userID string) *errors.Error
	DeleteAllInProject(projectName string) *errors.Error
}
