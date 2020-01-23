package model

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/role"
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

	// ErrUserValidateFailed ...
	ErrUserValidateFailed = errors.New("User validation failed")
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

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}

// Validate ...
func (ui *UserInfo) Validate() error {
	// Check User ID
	if !validateUserID(ui.ID) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid user ID format")
	}

	if !validateProjectName(ui.ProjectName) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid project Name format")
	}

	// Check User Name
	if !validateUserName(ui.Name) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid user name format")
	}

	// Check Roles
	for _, r := range ui.Roles {
		if !role.GetInst().IsValid(r) {
			return errors.Wrap(ErrUserValidateFailed, "Invalid role")
		}
	}

	return nil
}
