package db

import (
	"errors"
	"regexp"
	"time"
)

// RoleType ...
type RoleType int

const (
	// IDLengthMin is minimum length of ID
	IDLengthMin = 4
	// IDLengthMax is maximum length of ID
	IDLengthMax = 32
	// IDValidChar is regular expression of ID (only permit a-z, A-Z, 0-9, ., -, _)
	IDValidChar = `[^a-zA-Z0-9\.\-\_]`

	// PasswordLengthMin is minimum length of password
	PasswordLengthMin = 8
	// PasswordLengthMax is maximum length of password
	PasswordLengthMax = 128

	// RoleUserAdit allows user create/delete
	RoleUserAdit RoleType = 1
)

var (
	// ErrAuthFailed is an error for authentication failed
	ErrAuthFailed = errors.New("Failed to authenticate")
	// ErrUserAlreadyExists is an error for user is already exists
	ErrUserAlreadyExists = errors.New("User is already exists")
	// ErrNoSuchUser is an error for no such user
	ErrNoSuchUser = errors.New("No such user")
)

// User is a structure of user
type User struct {
	ID       string
	Password string
	Roles    []RoleType
}

//UserRequest is a request param for user method
type UserRequest struct {
	ID       string
	Password string
}

// TokenConfig is a structure for token config
type TokenConfig struct {
	ExpiredTime time.Time
	Issuer      string
}

// Validate method validates UserRequest
func (r *UserRequest) Validate() error {
	if len(r.ID) < IDLengthMin {
		return errors.New("ID Length is too small")
	}
	if len(r.ID) > IDLengthMax {
		return errors.New("ID Length is too long")
	}
	if regexp.MustCompile(IDValidChar).Match([]byte(r.ID)) {
		return errors.New("ID include unpermitted charactor")
	}

	if len(r.Password) < PasswordLengthMin {
		return errors.New("Password Length is too small")
	}
	if len(r.Password) > PasswordLengthMax {
		return errors.New("Password Length is too long")
	}

	return nil
}
