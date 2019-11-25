package model

import (
	"errors"
)

// UserInfo ...
type UserInfo struct {
	ID           string
	ProjectID    string
	Name         string
	Enabled      bool
	CreatedAt    string
	PasswordHash string
	Roles        []string
}

var (
	// ErrUserAlreadyExists ...
	ErrUserAlreadyExists = errors.New("User Already Exists")
)
