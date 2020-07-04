package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// CustomRole ...
type CustomRole struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	ProjectName string
}

// CustomRoleFilter ...
type CustomRoleFilter struct {
	Name string
}

var (
	// ErrNoSuchCustomRole ...
	ErrNoSuchCustomRole = errors.New("No Such Custom Role")

	// ErrCustomRoleAlreadyExists ...
	ErrCustomRoleAlreadyExists = errors.New("Custom Role Already Exists")

	// ErrCustomRoleValidateFailed ...
	ErrCustomRoleValidateFailed = errors.New("Custom Role Already Exists")
)

// CustomRoleHandler ...
type CustomRoleHandler interface {
	Add(projectName string, ent *CustomRole) *errors.Error
	Delete(projectName string, roleID string) *errors.Error
	Get(projectName string, roleID string) (*CustomRole, *errors.Error)
	GetList(projectName string, filter *CustomRoleFilter) ([]*CustomRole, *errors.Error)
	Update(projectName string, ent *CustomRole) *errors.Error
	DeleteAll(projectName string) *errors.Error
}

// Validate ...
func (c *CustomRole) Validate() *errors.Error {
	if !ValidateCustomRoleName(c.Name) {
		return errors.Append(ErrCustomRoleValidateFailed, "Invalid Custom Role Name format")
	}

	if !ValidateProjectName(c.ProjectName) {
		return errors.Append(ErrCustomRoleValidateFailed, "Invalid Project Name format")
	}

	return nil
}
