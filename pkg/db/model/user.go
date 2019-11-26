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

	// ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("No Such User")
)

// Validate ...
func (ui *UserInfo) Validate() error {
	if ui.ID == "" {
		return errors.New("User ID is empty") 
	}

	if ui.ProjectID == "" {
		return errors.New("Project ID is empty")
	}

	if ui.Name == "" {
		return errors.New("User Name is empty")
	}

	return nil
}