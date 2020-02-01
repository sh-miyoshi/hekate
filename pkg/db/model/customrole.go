package model

import (
	"github.com/pkg/errors"
	"time"
)

// CustomRole ...
type CustomRole struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	ProjectName string
}

var (
	// ErrNoSuchCustomRole ...
	ErrNoSuchCustomRole = errors.New("No Such Custom Role")

	// ErrCustomRoleAlreadyExists ...
	ErrCustomRoleAlreadyExists = errors.New("Custom Role Already Exists")
)

// CustomRoleHandler ...
type CustomRoleHandler interface {
	Add(ent *CustomRole) error
	Delete(roleID string) error
	Get(roleID string) (*CustomRole, error)
	GetList(projectName string) ([]string, error)
	Update(ent *CustomRole) error
	DeleteAll(projectName string) error

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}
