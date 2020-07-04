package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
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
	Add(projectName string, ent *UserInfo) *errors.Error
	Delete(projectName string, userID string) *errors.Error
	GetList(projectName string, filter *UserFilter) ([]*UserInfo, *errors.Error)
	Get(projectName string, userID string) (*UserInfo, *errors.Error)
	Update(projectName string, ent *UserInfo) *errors.Error
	DeleteAll(projectName string) *errors.Error
	AddRole(projectName string, userID string, roleType RoleType, roleID string) *errors.Error
	DeleteRole(projectName string, userID string, roleID string) *errors.Error
	DeleteAllCustomRole(projectName string, roleID string) *errors.Error
}

// Validate ...
func (ui *UserInfo) Validate() *errors.Error {
	// Check User ID
	if !ValidateUserID(ui.ID) {
		return errors.Append(ErrUserValidateFailed, "Invalid user ID format")
	}

	if !ValidateProjectName(ui.ProjectName) {
		return errors.Append(ErrUserValidateFailed, "Invalid project Name format")
	}

	// Check User Name
	if !ValidateUserName(ui.Name) {
		return errors.Append(ErrUserValidateFailed, "Invalid user name format")
	}

	return nil
}
