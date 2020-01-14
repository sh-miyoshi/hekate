package model

import (
	"errors"
	"time"
)

// AuthCode ...
type AuthCode struct {
	CodeID      string
	ExpiresIn   time.Time
	ClientID    string
	RedirectURL string
	UserID      string
}

// AuthCodeHandler ...
type AuthCodeHandler interface {
	New(code *AuthCode) error
	Get(codeID string) (*AuthCode, error)
	Delete(codeID string) error
}

var (
	// ErrCodeAlreadyExists ...
	ErrCodeAlreadyExists = errors.New("Code Already Exists")

	// ErrNoSuchCode ...
	ErrNoSuchCode = errors.New("No Such Code")
)
