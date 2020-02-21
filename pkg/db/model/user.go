package model

import (
	"time"

	"github.com/pkg/errors"
)

// UserInfo ...
type UserInfo struct {
	ID           string
	ProjectName  string
	Name         string
	CreatedAt    time.Time
	PasswordHash string
	SystemRoles  []string
	CustomRoles  []string
}

// UserFilter ...
type UserFilter struct {
	Name string
	// TODO(CreatedAt, SystemRoles, CustomRoles, ...)
}

// RoleType ...
type RoleType struct {
	value string
}

// String method returns a name of role type
func (t RoleType) String() string {
	return t.value
}

var (
	// ErrUserAlreadyExists ...
	ErrUserAlreadyExists = errors.New("User already exists")
	// ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("No such user")
	// ErrRoleAlreadyAppended ...
	ErrRoleAlreadyAppended = errors.New("Role already appended")
	// ErrNoSuchRoleInUser ...
	ErrNoSuchRoleInUser = errors.New("User do not have such role")
	// ErrUserValidateFailed ...
	ErrUserValidateFailed = errors.New("User validation failed")

	// RoleSystem ...
	RoleSystem = RoleType{"system_management"}
	// RoleCustom ...
	RoleCustom = RoleType{"custom_role"}
)

// UserInfoHandler ...
type UserInfoHandler interface {
	Add(ent *UserInfo) error
	Delete(userID string) error
	GetList(projectName string, filter *UserFilter) ([]*UserInfo, error)
	Get(userID string) (*UserInfo, error)
	Update(ent *UserInfo) error
	DeleteAll(projectName string) error
	AddRole(userID string, roleType RoleType, roleID string) error
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
	if !ValidateUserID(ui.ID) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid user ID format")
	}

	if !ValidateProjectName(ui.ProjectName) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid project Name format")
	}

	// Check User Name
	if !ValidateUserName(ui.Name) {
		return errors.Wrap(ErrUserValidateFailed, "Invalid user name format")
	}

	return nil
}
