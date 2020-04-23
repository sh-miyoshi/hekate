package model

import (
	"time"

	"github.com/pkg/errors"
)

// LoginSessionInfo ...
type LoginSessionInfo struct {
	VerifyCode   string
	ExpiresIn    time.Time
	Scope        string
	ResponseType string
	ClientID     string
	RedirectURI  string
	Nonce        string
	ProjectName  string
}

var (
	// ErrLoginSessionAlreadyExists ...
	ErrLoginSessionAlreadyExists = errors.New("Login session already exists")
	// ErrNoSuchLoginSession ...
	ErrNoSuchLoginSession = errors.New("No such login session")
	// ErrLoginSessionValidationFailed ...
	ErrLoginSessionValidationFailed = errors.New("login session validation failed")
)

// LoginSessionHandler ...
type LoginSessionHandler interface {
	Add(info *LoginSessionInfo) error
	Delete(code string) error
	Get(code string) (*LoginSessionInfo, error)
	DeleteAll(clientID string) error
	DeleteAllInProject(projectName string) error
}
