package model

import (
	"github.com/pkg/errors"
	"time"
	"github.com/asaskevich/govalidator"
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

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}

var (
	// ErrCodeAlreadyExists ...
	ErrCodeAlreadyExists = errors.New("Code Already Exists")

	// ErrNoSuchCode ...
	ErrNoSuchCode = errors.New("No Such Code")

	// ErrCodeValidateFailed ...
	ErrCodeValidateFailed = errors.New("Code validation failed")
)

// Validate ...
func (ac *AuthCode) Validate() error {
	if !govalidator.IsUUID(ac.CodeID) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid code ID format")
	}

	if !validateClientID(ac.ClientID) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid client ID format")
	}

	if !validateUserID(ac.UserID) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid user ID format")
	}

	if !govalidator.IsURL(ac.RedirectURL) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid redirect url format")
	}

	return nil
}