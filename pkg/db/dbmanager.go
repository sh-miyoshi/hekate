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
	session model.SessionHandler
	client  model.ClientInfoHandler
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
		sessionHander, err := memory.NewSessionHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create session handler")
		}
		clientHandler, err := memory.NewClientHandler(prjHandler)
		if err != nil {
			return errors.Wrap(err, "Failed to create client handler")
		}

		inst = &Manager{
			project: prjHandler,
			user:    userHandler,
			session: sessionHander,
			client:  clientHandler,
		}
	case "mongo":
		logger.Info("Initialize with mongo DB")
		dbClient, err := mongo.NewClient(connStr)
		if err != nil {
			return errors.Wrap(err, "Failed to create db client")
		}

		prjHandler, err := mongo.NewProjectHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create project handler")
		}
		userHandler, err := mongo.NewUserHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create user handler")
		}
		sessionHandler, err := mongo.NewSessionHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create session handler")
		}
		clientHandler, err := mongo.NewClientHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create client handler")
		}

		inst = &Manager{
			project: prjHandler,
			user:    userHandler,
			session: sessionHandler,
			client:  clientHandler,
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

	if err := m.user.DeleteAll(name); err != nil {
		return errors.Wrap(err, "failed to delete user data")
	}

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

	_, err := m.user.Get(ent.ID)
	if err != model.ErrNoSuchUser {
		if err == nil {
			return errors.Cause(model.ErrUserAlreadyExists)
		}
		return errors.Wrap(err, "Failed to get user info")
	}

	// Check duplicate user by name
	_, err = m.user.GetByName(ent.ProjectName, ent.Name)
	if err != model.ErrNoSuchUser {
		if err == nil {
			return errors.Cause(model.ErrUserAlreadyExists)
		}
		return errors.Wrap(err, "Failed to get user info by name")
	}

	return m.user.Add(ent)
}

// UserDelete ...
func (m *Manager) UserDelete(userID string) error {
	// TODO(validate userID)
	// TODO(lock, unlock)
	return m.user.Delete(userID)
}

// UserGetList ...
func (m *Manager) UserGetList(projectName string) ([]string, error) {
	// TODO(validate projectName)
	return m.user.GetList(projectName)
}

// UserGet ...
func (m *Manager) UserGet(userID string) (*model.UserInfo, error) {
	// TODO(validate userID)
	return m.user.Get(userID)
}

// UserUpdate ...
func (m *Manager) UserUpdate(ent *model.UserInfo) error {
	// TODO(validate ent)
	// TODO(lock, unlock)
	return m.user.Update(ent)
}

// UserGetByName ...
func (m *Manager) UserGetByName(projectName string, userName string) (*model.UserInfo, error) {
	// TODO(validate projectName, userName)
	return m.user.GetByName(projectName, userName)
}

// UserAddRole ...
func (m *Manager) UserAddRole(userID string, roleID string) error {
	// TODO(validate userID, roleID)
	// TODO(lock, unlock)
	return m.user.AddRole(userID, roleID)
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(userID string, roleID string) error {
	// TODO(validate userID, roleID)
	// TODO(lock, unlock)
	return m.user.DeleteRole(userID, roleID)
}

// NewSession ...
func (m *Manager) NewSession(userID string, sessionID string, expiresIn uint, fromIP string) error {
	// TODO(validate userID, sessionID, expiresIn?)
	// TODO(lock, unlock)

	if _, err := m.session.Get(sessionID); err != model.ErrNoSuchSession {
		return errors.Cause(model.ErrSessionAlreadyExists)
	}

	return m.session.New(userID, sessionID, expiresIn, fromIP)
}

// RevokeSession ...
func (m *Manager) RevokeSession(sessionID string) error {
	// TODO(validate sessionID)
	// TODO(lock, unlock)
	return m.session.Revoke(sessionID)
}

// GetSessions ...
func (m *Manager) GetSessions(userID string) ([]string, error) {
	// TODO(validate userID)
	return m.session.GetList(userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(ent *model.ClientInfo) error {
	// TODO(validation)
	// TODO(lock, unlock)

	_, err := m.client.Get(ent.ID)
	if err != model.ErrNoSuchClient {
		if err == nil {
			return errors.Cause(model.ErrClientAlreadyExists)
		}
		return errors.Wrap(err, "Failed to get client info")
	}

	return m.client.Add(ent)
}

// ClientDelete ...
func (m *Manager) ClientDelete(clientID string) error {
	// TODO(validate clientID)
	// TODO(lock, unlock)
	return m.client.Delete(clientID)
}

// ClientGetList ...
func (m *Manager) ClientGetList(projectName string) ([]string, error) {
	// TODO(validate projectName)
	return m.client.GetList(projectName)
}

// ClientGet ...
func (m *Manager) ClientGet(clientID string) (*model.ClientInfo, error) {
	// TODO(validate clientID)
	return m.client.Get(clientID)
}

// ClientUpdate ...
func (m *Manager) ClientUpdate(ent *model.ClientInfo) error {
	// TODO(validate ent)
	// TODO(lock, unlock)
	return m.client.Update(ent)
}
