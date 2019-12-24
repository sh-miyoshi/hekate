package model

import (
	"errors"
	"time"
)

// UserInfo ...
type UserInfo struct {
	ID           string
	ProjectName  string
	Name         string
	CreatedAt    time.Time
	PasswordHash string
	Roles        []string
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

// UserInfoHandler ...
type UserInfoHandler interface {
	Add(ent *UserInfo) error
	Delete(userID string) error
	GetList(projectName string) ([]string, error)
	Get(userID string) (*UserInfo, error)
	GetByName(projectName string, userName string) (*UserInfo, error)
	Update(ent *UserInfo) error
	DeleteAll(projectName string) error
	AddRole(userID string, roleID string) error
	DeleteRole(userID string, roleID string) error
}

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
