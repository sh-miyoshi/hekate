package model

import (
	"github.com/pkg/errors"
	"time"
)

// LoginSessionInfo ...
type LoginSessionInfo struct {
	VerifyCode  string
	ExpiresIn   time.Time
	ClientID    string
	RedirectURI string
}

var (
	// ErrLoginSessionAlreadyExists ...
	ErrLoginSessionAlreadyExists = errors.New("Login session already exists")
	// ErrNoSuchLoginSession ...
	ErrNoSuchLoginSession = errors.New("No such login session")
)

// LoginSessionHandler ...
type LoginSessionHandler interface {
	Add(info *LoginSessionInfo) error
	Delete(code string) error
	Get(code string) (*LoginSessionInfo, error)

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}
