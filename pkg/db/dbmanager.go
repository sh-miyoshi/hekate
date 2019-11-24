package db

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/db/local"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Manager ...
type Manager struct {
	Project ProjectInfoHandler
}

var inst *Manager

// InitDBManager ...
func InitDBManager(dbType string, connStr string) error {
	if inst != nil {
		return fmt.Errorf("DBManager is already initialized")
	}

	switch dbType {
	case "local":
		logger.Info("Initialize with local file storage")
		prjHandler, err := local.NewHandler(connStr)
		if err != nil {
			return err
		}
		inst = &Manager{
			Project: prjHandler,
		}
	default:
		return fmt.Errorf("Database Type %s is not implemented yet", dbType)
	}

	return nil
}
