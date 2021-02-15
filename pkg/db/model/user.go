package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// LockState ...
type LockState struct {
	Locked            bool
	VerifyFailedTimes []time.Time
}

// OTPInfo ...
type OTPInfo struct {
	ID         string
	PrivateKey string
	Enabled    bool
}

// UserInfo ...
type UserInfo struct {
	ID           string
	ProjectName  string
	Name         string
	CreatedAt    time.Time
	PasswordHash string
	SystemRoles  []string
	CustomRoles  []string
	LockState    LockState
	OTPInfo      OTPInfo
}

// UserFilter ...
type UserFilter struct {
	ID   string
	Name string
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
	ErrUserAlreadyExists = errors.New("User already exists", "User already exists")
	// ErrNoSuchUser ...
	ErrNoSuchUser = errors.New("No such user", "No such user")
	// ErrRoleAlreadyAppended ...
	ErrRoleAlreadyAppended = errors.New("Role already appended", "Role already appended")
	// ErrNoSuchRoleInUser ...
	ErrNoSuchRoleInUser = errors.New("User do not have such role", "User do not have such role")
	// ErrUserValidateFailed ...
	ErrUserValidateFailed = errors.New("User validation failed", "User validation failed")
	// ErrUserOTPAlreadyEnabled ...
	ErrUserOTPAlreadyEnabled = errors.New("User OTP already enabled", "User OTP already enabled")

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
