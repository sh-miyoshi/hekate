package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/memory"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/db/mongo"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Manager ...
type Manager struct {
	project model.ProjectInfoHandler
	user    model.UserInfoHandler
}

var inst *Manager

// InitDBManager ...
func InitDBManager(dbType string, connStr string) error {
	if inst != nil {
		return errors.Cause(fmt.Errorf("DBManager is already initialized"))
	}

	switch dbType {
	case "memory":
		logger.Info("Initialize with local memory DB")
		prjHandler, err := memory.NewProjectHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create project handler")
		}
		userHandler, err := memory.NewUserHandler(prjHandler)
		if err != nil {
			return errors.Wrap(err, "Failed to create user handler")
		}

		inst = &Manager{
			project: prjHandler,
			user:    userHandler,
		}
	case "mongo":
		logger.Info("Initialize with mongo DB")
		dbClient, err := mongo.NewClient(connStr)
		if err != nil {
			return errors.Wrap(err, "Failed to create db client")
		}

		prjHandler := mongo.NewProjectHandler(dbClient)
		userHandler := mongo.NewUserHandler(dbClient, prjHandler)

		inst = &Manager{
			project: prjHandler,
			user:    userHandler,
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

// ProjectAdd ...
func (m *Manager) ProjectAdd(ent *model.ProjectInfo) error {
	if ent.Name == "" {
		return errors.New("name of entry is empty")
	}

	// TODO(lock, unlock)

	if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
		return errors.Cause(model.ErrProjectAlreadyExists)
	}

	return m.project.Add(ent)
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) error {
	if name == "" {
		return errors.New("name of entry is empty")
	}

	if name == "master" {
		return errors.Wrap(model.ErrDeleteBlockedProject, "master project can not delete")
	}

	// TODO(lock, unlock)

	return m.project.Delete(name)
}

// ProjectGetList ...
func (m *Manager) ProjectGetList() ([]string, error) {
	return m.project.GetList()
}

// ProjectGet ...
func (m *Manager) ProjectGet(name string) (*model.ProjectInfo, error) {
	if name == "" {
		return nil, errors.New("name of entry is empty")
	}

	return m.project.Get(name)
}

// ProjectUpdate ...
func (m *Manager) ProjectUpdate(ent *model.ProjectInfo) error {
	// TODO(validate ent, projectExists)
	// TODO(lock, unlock)
	return m.project.Update(ent)
}

// UserAdd ...
func (m *Manager) UserAdd(ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// TODO(lock, unlock)

	_, err := m.user.Get(ent.ProjectName, ent.ID)
	if err != model.ErrNoSuchUser {
		if err == nil {
			return errors.Cause(model.ErrUserAlreadyExists)
		}
		return errors.Wrap(err, "Failed to get user info")
	}

	// Check duplicate user by name
	if id, err := m.user.GetIDByName(ent.ProjectName, ent.Name); err == nil {
		_, err = m.user.Get(ent.ProjectName, id)
		if err != model.ErrNoSuchUser {
			if err == nil {
				return errors.Cause(model.ErrUserAlreadyExists)
			}
			return errors.Wrap(err, "Failed to get user info by name")
		}
	}

	return m.user.Add(ent)
}

// UserDelete ...
func (m *Manager) UserDelete(projectName string, userID string) error {
	// TODO(validate projectName, userID)
	// TODO(lock, unlock)
	return m.user.Delete(projectName, userID)
}

// UserGetList ...
func (m *Manager) UserGetList(projectName string) ([]string, error) {
	// TODO(validate projectName)
	return m.user.GetList(projectName)
}

// UserGet ...
func (m *Manager) UserGet(projectName string, userID string) (*model.UserInfo, error) {
	// TODO(validate projectName, userID)
	return m.user.Get(projectName, userID)
}

// UserUpdate ...
func (m *Manager) UserUpdate(ent *model.UserInfo) error {
	// TODO(validate ent)
	// TODO(lock, unlock)
	return m.user.Update(ent)
}

// UserGetIDByName ...
func (m *Manager) UserGetIDByName(projectName string, userName string) (string, error) {
	// TODO(validate projectName, userName)
	return m.user.GetIDByName(projectName, userName)
}

// UserDeleteAll ...
func (m *Manager) UserDeleteAll(projectName string) error {
	// TODO(validate projectName)
	// TODO(lock, unlock)
	return m.user.DeleteAll(projectName)
}

// UserAddRole ...
func (m *Manager) UserAddRole(projectName string, userID string, roleID string) error {
	// TODO(validate projectName, userID, roleID)
	// TODO(lock, unlock)
	return m.user.AddRole(projectName, userID, roleID)
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(projectName string, userID string, roleID string) error {
	// TODO(validate projectName, userID, roleID)
	// TODO(lock, unlock)
	return m.user.DeleteRole(projectName, userID, roleID)
}

// NewSession ...
func (m *Manager) NewSession(projectName string, userID string, session model.Session) error {
	// TODO(validate projectName, userID, session)
	// TODO(lock, unlock)
	return m.user.NewSession(projectName, userID, session)
}

// RevokeSession ...
func (m *Manager) RevokeSession(projectName string, userID string, sessionID string) error {
	// TODO(validate projectName, userID, sessionID)
	// TODO(lock, unlock)
	return m.user.RevokeSession(projectName, userID, sessionID)
}
