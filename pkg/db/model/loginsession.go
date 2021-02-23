package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// LoginSession ...
type LoginSession struct {
	SessionID           string
	Code                string
	ExpiresDate         time.Time
	Scope               string
	ResponseType        []string
	ClientID            string
	RedirectURI         string
	Nonce               string
	ProjectName         string
	ResponseMode        string
	Prompt              []string
	UserID              string
	LoginDate           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
}

// LoginSessionFilter ...
type LoginSessionFilter struct {
	SessionID string
	UserID    string
	ClientID  string
}

var (
	// ErrNoSuchLoginSession ...
	ErrNoSuchLoginSession = errors.New("No such session", "No such session")
	// ErrLoginSessionValidationFailed ...
	ErrLoginSessionValidationFailed = errors.New("Login Session validation failed", "Login Session validation failed")
)

// LoginSessionHandler ...
type LoginSessionHandler interface {
	Add(projectName string, ent *LoginSession) *errors.Error
	Update(projectName string, ent *LoginSession) *errors.Error
	Delete(projectName string, filter *LoginSessionFilter) *errors.Error
	DeleteAll(projectName string) *errors.Error
	GetByCode(projectName string, code string) (*LoginSession, *errors.Error)
	Get(projectName string, sessionID string) (*LoginSession, *errors.Error)
	Cleanup(now time.Time) *errors.Error
}
