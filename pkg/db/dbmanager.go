package db

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/memory"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/db/mongo"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Manager ...
type Manager struct {
	project  model.ProjectInfoHandler
	user     model.UserInfoHandler
	session  model.SessionHandler
	client   model.ClientInfoHandler
	authCode model.AuthCodeHandler
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
		authCodeHandler, err := memory.NewAuthCodeHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create auth code handler")
		}

		inst = &Manager{
			project:  prjHandler,
			user:     userHandler,
			session:  sessionHander,
			client:   clientHandler,
			authCode: authCodeHandler,
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
		authCodeHandler, err := mongo.NewAuthCodeHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create auth code handler")
		}

		inst = &Manager{
			project:  prjHandler,
			user:     userHandler,
			session:  sessionHandler,
			client:   clientHandler,
			authCode: authCodeHandler,
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
	if err := ent.Validate(); err != nil {
		logger.Info("Failed to validate project entry: %v", err)
		return errors.Cause(model.ErrProjectValidationFailed)
	}

	switch ent.TokenConfig.SigningAlgorithm {
	case "RS256":
		key, err := rsa.GenerateKey(rand.Reader, 2048) // fixed key length is ok?
		if err != nil {
			return errors.Wrap(err, "Failed to generate RSA private key")
		}
		ent.TokenConfig.SignSecretKey = x509.MarshalPKCS1PrivateKey(key)
		ent.TokenConfig.SignPublicKey = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	}

	if err := m.project.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
		m.project.AbortTx()
		return errors.Cause(model.ErrProjectAlreadyExists)
	}

	if err := m.project.Add(ent); err != nil {
		m.project.AbortTx()
		return err
	}
	m.project.CommitTx()
	return nil
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) error {
	if name == "" {
		return errors.New("name of entry is empty")
	}

	if name == "master" {
		return errors.Wrap(model.ErrDeleteBlockedProject, "master project can not delete")
	}

	if err := m.project.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	if err := m.user.DeleteAll(name); err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "failed to delete user data")
	}

	if err := m.project.Delete(name); err != nil {
		m.project.AbortTx()
		return err
	}
	m.project.CommitTx()
	return nil
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
	if err := ent.Validate(); err != nil {
		logger.Info("Failed to validate project entry: %v", err)
		return errors.Cause(model.ErrProjectValidationFailed)
	}

	if err := m.project.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	if err := m.project.Update(ent); err != nil {
		m.project.AbortTx()
		return err
	}
	m.project.CommitTx()
	return nil
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
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

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

// SessionAdd ...
func (m *Manager) SessionAdd(ent *model.Session) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}
	// TODO(lock, unlock)

	if _, err := m.session.Get(ent.SessionID); err != model.ErrNoSuchSession {
		return errors.Cause(model.ErrSessionAlreadyExists)
	}

	return m.session.New(ent)
}

// SessionDelete ...
func (m *Manager) SessionDelete(sessionID string) error {
	// TODO(validate sessionID)
	// TODO(lock, unlock)
	return m.session.Revoke(sessionID)
}

// SessionGetList ...
func (m *Manager) SessionGetList(userID string) ([]string, error) {
	// TODO(validate userID)
	return m.session.GetList(userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if err := m.client.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	_, err := m.client.Get(ent.ID)
	if err != model.ErrNoSuchClient {
		m.client.AbortTx()
		if err == nil {
			return errors.Cause(model.ErrClientAlreadyExists)
		}
		return errors.Wrap(err, "Failed to get client info")
	}

	if err := m.client.Add(ent); err != nil {
		m.client.AbortTx()
		return err
	}
	m.client.CommitTx()
	return nil
}

// ClientDelete ...
func (m *Manager) ClientDelete(clientID string) error {
	// TODO(validate clientID)

	if err := m.client.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	if err := m.client.Delete(clientID); err != nil {
		m.client.AbortTx()
		return err
	}
	m.client.CommitTx()
	return nil
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
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if err := m.client.BeginTx(); err != nil {
		return errors.Cause(err)
	}

	if err := m.client.Update(ent); err != nil {
		m.client.AbortTx()
		return err
	}
	m.client.CommitTx()
	return nil
}

// AuthCodeAdd ...
func (m *Manager) AuthCodeAdd(ent *model.AuthCode) error {
	// TODO(validate ent, identify by clientID and redirectURL)
	// TODO(lock, unlock)
	return m.authCode.New(ent)
}

// AuthCodeDelete ...
func (m *Manager) AuthCodeDelete(codeID string) error {
	// TODO(validate codeID)
	// TODO(lock, unlock)
	return m.authCode.Delete(codeID)
}

// AuthCodeGet ...
func (m *Manager) AuthCodeGet(codeID string) (*model.AuthCode, error) {
	// TODO(validate codeID)
	return m.authCode.Get(codeID)
}
