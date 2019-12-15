package model

import (
	"errors"
	"time"
)

// Session ...
type Session struct {
	SessionID string
	CreatedAt time.Time
	ExpiresIn uint
	FromIP    string // Used to identify the user using this session
}

// UserInfo ...
type UserInfo struct {
	ID           string
	ProjectName  string
	Name         string
	CreatedAt    time.Time
	PasswordHash string
	Roles        []string
	Sessions     []Session
}

var (
	// ErrUserAlreadyExists ...
	ErrUserAlreadyExists = errors.New("User Already Exists")

	// ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("No Such User")

	// ErrRoleAlreadyAppended ...
	ErrRoleAlreadyAppended = errors.New("Role already appended")

	// ErrNoSuchRoleInUser ...
	ErrNoSuchRoleInUser = errors.New("User do not have such role")
)

// Validate ...
func (ui *UserInfo) Validate() error {
	if ui.ID == "" {
		return errors.New("User ID is empty")
	}

	if ui.ProjectName == "" {
		return errors.New("Project Name is empty")
	}

	if ui.Name == "" {
		return errors.New("User Name is empty")
	}

	return nil
}
