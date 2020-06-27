package model

import (
	"time"

	"github.com/pkg/errors"
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
	Add(projectName string, ent *CustomRole) error
	Delete(projectName string, roleID string) error
	Get(projectName string, roleID string) (*CustomRole, error)
	GetList(projectName string, filter *CustomRoleFilter) ([]*CustomRole, error)
	Update(projectName string, ent *CustomRole) error
	DeleteAll(projectName string) error
}

// Validate ...
func (c *CustomRole) Validate() error {
	if !ValidateCustomRoleName(c.Name) {
		return errors.Wrap(ErrCustomRoleValidateFailed, "Invalid Custom Role Name format")
	}

	if !ValidateProjectName(c.ProjectName) {
		return errors.Wrap(ErrCustomRoleValidateFailed, "Invalid Project Name format")
	}

	return nil
}
