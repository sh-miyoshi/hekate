package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

// AuthCode ...
type AuthCode struct {
	CodeID      string
	ExpiresIn   time.Time
	ClientID    string
	RedirectURL string
	UserID      string
	Nonce       string
	MaxAge      uint
}

// AuthCodeHandler ...
type AuthCodeHandler interface {
	Add(code *AuthCode) error
	Get(codeID string) (*AuthCode, error)
	Delete(codeID string) error
	DeleteAll(userID string) error
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

	if !ValidateClientID(ac.ClientID) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid client ID format")
	}

	if !ValidateUserID(ac.UserID) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid user ID format")
	}

	if !govalidator.IsURL(ac.RedirectURL) {
		return errors.Wrap(ErrCodeValidateFailed, "Invalid redirect url format")
	}

	return nil
}
