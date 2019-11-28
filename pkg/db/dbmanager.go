package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/local"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Manager ...
type Manager struct {
	Project ProjectInfoHandler
	User    UserInfoHandler
}

var inst *Manager

// InitDBManager ...
func InitDBManager(dbType string, connStr string) error {
	if inst != nil {
		return errors.Cause(fmt.Errorf("DBManager is already initialized"))
	}

	switch dbType {
	case "local":
		logger.Info("Initialize with local file storage")
		prjHandler, err := local.NewProjectHandler(connStr)
		if err != nil {
			return errors.Wrap(err, "Failed to create project handler")
		}
		userHandler, err := local.NewUserHandler(connStr)
		if err != nil {
			return errors.Wrap(err, "Failed to create user handler")
		}

		inst = &Manager{
			Project: prjHandler,
			User:    userHandler,
		}
	default:
		return errors.Cause(fmt.Errorf("Database Type %s is not implemented yet", dbType))
	}

	return nil
}

// GetInst returns an instance of DB Manager
func GetInst() *Manager {
	return inst
}
