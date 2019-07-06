package db

import (
	"errors"
	"regexp"
	"time"
)

// RoleType ...
type RoleType int

const (
	// NameLengthMin is minimum length of name
	NameLengthMin = 4
	// NameLengthMax is maximum length of name
	NameLengthMax = 32
	// NameValidChar is regular expression of name (only permit a-z, A-Z, 0-9, ., -, _)
	NameValidChar = `[^a-zA-Z0-9\.\-\_]`

	// PasswordLengthMin is minimum length of password
	PasswordLengthMin = 8
	// PasswordLengthMax is maximum length of password
	PasswordLengthMax = 128

	// RoleAdmin is admin role
	RoleAdmin RoleType = 1
	// RoleWrite is write role
	RoleWrite RoleType = 2
	// RoleRead is read role
	RoleRead RoleType = 4
)

var (
	// ErrAuthFailed is an error for authentication failed
	ErrAuthFailed = errors.New("Failed to authenticate")
	// ErrUserAlreadyExists is an error for user is already exists
	ErrUserAlreadyExists = errors.New("User is already exists")
	// ErrNoSuchUser is an error for no such user
	ErrNoSuchUser = errors.New("No such user")
)

// Role ...
type Role struct {
	ID   string
	Type RoleType
	// TargetResource TODO(now only user)
}

// User is a structure of user
type User struct {
	ID       string
	Name     string
	Password string
	Roles    []Role
}

//UserRequest is a request param for user method
type UserRequest struct {
	Name     string
	Password string
}

// TokenConfig is a structure for token config
type TokenConfig struct {
	ExpiredTime time.Time
	Issuer      string
}

// Validate method validates UserRequest
func (r *UserRequest) Validate() error {
	if len(r.Name) < NameLengthMin {
		return errors.New("Name Length is too small")
	}
	if len(r.Name) > NameLengthMax {
		return errors.New("Name Length is too long")
	}
	if regexp.MustCompile(NameValidChar).Match([]byte(r.Name)) {
		return errors.New("Name include unpermitted charactor")
	}

	if len(r.Password) < PasswordLengthMin {
		return errors.New("Password Length is too small")
	}
	if len(r.Password) > PasswordLengthMax {
		return errors.New("Password Length is too long")
	}

	return nil
}
