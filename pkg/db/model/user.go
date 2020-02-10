package model

import (
	"github.com/pkg/errors"
	"time"
)

// LoginSessionInfo ...
type LoginSessionInfo struct {
	VerifyCode  string
	ExpiresIn   time.Time
	ClientID    string
	RedirectURI string
}

// UserInfo ...
type UserInfo struct {
	ID            string
	ProjectName   string
	Name          string
	CreatedAt     time.Time
	PasswordHash  string
	SystemRoles   []string
	CustomRoles   []string
	LoginSessions []*LoginSessionInfo
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
	// ErrLoginSessionAlreadyExists ...
	ErrLoginSessionAlreadyExists = errors.New("Login session already exists")
	// ErrNoSuchLoginSession ...
	ErrNoSuchLoginSession = errors.New("No such login session")

	// RoleSystem ...
	RoleSystem = RoleType{"system_management"}
	// RoleCustom ...
	RoleCustom = RoleType{"custom_role"}
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
	AddRole(userID string, roleType RoleType, roleID string) error
	DeleteRole(userID string, roleID string) error
	AddLoginSession(userID string, info *LoginSessionInfo) error
	DeleteLoginSession(userID string, code string) error

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
