package audit

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/memory"
	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/audit/mongo"
	"github.com/sh-miyoshi/hekate/pkg/audit/none"
	dbmongo "github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

// Manager ...
type Manager struct {
	handler Handler
}

var inst *Manager

// Init ...
func Init(dbType string, connStr string) *errors.Error {
	if inst != nil {
		return errors.New("Internal server error", "AuditManager is already initialized")
	}

	switch dbType {
	case "memory":
		logger.Info("Initialize AuditManager with local memory DB")

		inst = &Manager{
			handler: memory.NewHandler(),
		}
	case "mongo":
		logger.Info("Initialize AuditManager with mongo DB")
		dbClient, err := dbmongo.NewClient(connStr)
		if err != nil {
			return errors.Append(err, "Failed to get db client")
		}
		inst = &Manager{
			handler: mongo.NewHandler(dbClient),
		}
	case "none":
		logger.Info("Initialize AuditManager with none DB")
		inst = &Manager{
			handler: none.NewHandler(),
		}
	default:
		return errors.New("Internal server error", "Database Type %s is not implemented for audit events", dbType)
	}
	return nil
}

// GetInst returns an instance of DB Manager
func GetInst() *Manager {
	return inst
}

// Ping ...
func (m *Manager) Ping() *errors.Error {
	return m.handler.Ping()
}

// Save ...
func (m *Manager) Save(projectName string, tm time.Time, resType, method, path, message string) *errors.Error {
	return m.handler.Save(projectName, tm, resType, method, path, message)
}

// Get ...
func (m *Manager) Get(projectName string, fromDate, toDate time.Time, offset uint) ([]model.Audit, *errors.Error) {
	fromDate = util.TimeTruncate(fromDate)
	toDate = util.TimeTruncate(toDate)
	return m.handler.Get(projectName, fromDate, toDate, offset)
}
