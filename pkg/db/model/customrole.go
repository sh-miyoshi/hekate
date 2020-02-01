package model

import (
	"time"
)

// CustomRole ...
type CustomRole struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	ProjectName string
}

// CustomRoleHandler ...
type CustomRoleHandler interface {
	Add(ent *CustomRole) error
	Delete(roleID string) error
	Get(roleID string) (*CustomRole, error)
	Update(ent *CustomRole) error
	DeleteAll(projectName string) error

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}
