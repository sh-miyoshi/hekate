package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/role"
	"regexp"
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
	// Check User ID
	if ok := govalidator.IsUUID(ui.ID); !ok {
		return errors.New("Invalid user ID format")
	}

	// Check Project Name
	prjNameRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-]{2,31}$`)
	if !prjNameRegExp.MatchString(ui.ProjectName) {
		return errors.New("Invalid project Name format")
	}

	// Check User Name
	if !(3 <= len(ui.Name) && len(ui.Name) < 64) {
		return errors.New("Invalid user name format")
	}

	// Check Roles
	for _, r := range ui.Roles {
		if !role.GetInst().IsValid(r) {
			return errors.New("Invalid role")
		}
	}

	return nil
}
