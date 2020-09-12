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
	ID   string
	Name string
}

var (
	// ErrNoSuchCustomRole ...
	ErrNoSuchCustomRole = errors.New("No such custom role", "No such custom role")

	// ErrCustomRoleAlreadyExists ...
	ErrCustomRoleAlreadyExists = errors.New("Custom role already exists", "Custom role already exists")

	// ErrCustomRoleValidateFailed ...
	ErrCustomRoleValidateFailed = errors.New("Custom role validation failed", "Custom role validation failed")
)

// CustomRoleHandler ...
type CustomRoleHandler interface {
	Add(projectName string, ent *CustomRole) *errors.Error
	Delete(projectName string, roleID string) *errors.Error
	GetList(projectName string, filter *CustomRoleFilter) ([]*CustomRole, *errors.Error)
	Update(projectName string, ent *CustomRole) *errors.Error
	DeleteAll(projectName string) *errors.Error
}

// Validate ...
func (c *CustomRole) Validate() *errors.Error {
	if !ValidateCustomRoleID(c.ID) {
		return errors.Append(ErrCustomRoleValidateFailed, "Invalid Custom Role ID format")
	}

	if !ValidateCustomRoleName(c.Name) {
		return errors.Append(ErrCustomRoleValidateFailed, "Invalid Custom Role Name format")
	}

	if !ValidateProjectName(c.ProjectName) {
		return errors.Append(ErrCustomRoleValidateFailed, "Invalid Project Name format")
	}

	return nil
}
