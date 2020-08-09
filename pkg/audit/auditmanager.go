package audit

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/memory"
	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// Manager ...
type Manager struct {
	handler Handler
}

var inst *Manager

// Init ...
func Init(dbType string, connStr string) *errors.Error {
	if inst != nil {
		return errors.New("", "AuditManager is already initialized")
	}

	switch dbType {
	case "memory":
		logger.Info("Initialize AuditManager with local memory DB")

		inst = &Manager{
			handler: memory.NewHandler(),
		}
	default:
		return errors.New("", "Database Type %s is not implemented for Audit logs", dbType)
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
func (m *Manager) Save(tm time.Time, resType, method, path, body string) *errors.Error {
	return m.handler.Save(tm, resType, method, path, body)
}

// Get ...
func (m *Manager) Get(fromDate, toDate time.Time) ([]model.Audit, *errors.Error) {
	return m.handler.Get(fromDate, toDate)
}
