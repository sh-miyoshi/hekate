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
	"github.com/sh-miyoshi/jwt-server/pkg/role"
)

// Manager ...
type Manager struct {
	project      model.ProjectInfoHandler
	user         model.UserInfoHandler
	session      model.SessionHandler
	client       model.ClientInfoHandler
	authCode     model.AuthCodeHandler
	customRole   model.CustomRoleHandler
	loginSession model.LoginSessionHandler
}

var inst *Manager

// InitDBManager ...
func InitDBManager(dbType string, connStr string) error {
	if inst != nil {
		return errors.New(fmt.Sprintf("DBManager is already initialized"))
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
		clientHandler, err := memory.NewClientHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create client handler")
		}
		authCodeHandler, err := memory.NewAuthCodeHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create auth code handler")
		}
		customRoleHandler, err := memory.NewCustomRoleHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create custom role handler")
		}
		loginSessionHandler, err := memory.NewLoginSessionHandler()
		if err != nil {
			return errors.Wrap(err, "Failed to create login session handler")
		}

		inst = &Manager{
			project:      prjHandler,
			user:         userHandler,
			session:      sessionHander,
			client:       clientHandler,
			authCode:     authCodeHandler,
			customRole:   customRoleHandler,
			loginSession: loginSessionHandler,
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
		customRoleHandler, err := mongo.NewCustomRoleHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create custom role handler")
		}
		loginSessionHandler, err := mongo.NewLoginSessionHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create login session handler")
		}

		inst = &Manager{
			project:      prjHandler,
			user:         userHandler,
			session:      sessionHandler,
			client:       clientHandler,
			authCode:     authCodeHandler,
			customRole:   customRoleHandler,
			loginSession: loginSessionHandler,
		}
	default:
		return errors.New(fmt.Sprintf("Database Type %s is not implemented yet", dbType))
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
		return errors.Wrap(err, "Validate failed")
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
		return errors.Wrap(err, "BeginTx failed")
	}

	if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
		m.project.AbortTx()
		return model.ErrProjectAlreadyExists
	}

	if err := m.project.Add(ent); err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "Failed to add project")
	}
	m.project.CommitTx()
	return nil
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) error {
	if name == "" {
		return errors.New("name of entry is empty")
	}

	if err := m.project.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	prj, err := m.project.Get(name)
	if err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "Failed to get delete project info")
	}

	if !prj.PermitDelete {
		m.project.AbortTx()
		return errors.Wrap(model.ErrDeleteBlockedProject, "the project can not delete")
	}

	if err := m.user.DeleteAll(name); err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "Failed to delete user data")
	}

	if err := m.project.Delete(name); err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "Failed to delete project")
	}
	m.project.CommitTx()
	return nil
}

// ProjectGetList ...
func (m *Manager) ProjectGetList() ([]*model.ProjectInfo, error) {
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
		return errors.Wrap(err, "Failed to validate")
	}

	if err := m.project.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.project.Update(ent); err != nil {
		m.project.AbortTx()
		return errors.Wrap(err, "Failed to update project")
	}
	m.project.CommitTx()
	return nil
}

// UserAdd ...
func (m *Manager) UserAdd(ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// Validate Roles
	for _, r := range ent.SystemRoles {
		if !role.GetInst().IsValid(r) {
			return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
		}
	}
	for _, r := range ent.CustomRoles {
		if _, err := m.customRole.Get(r); err != nil {
			if errors.Cause(err) == model.ErrNoSuchCustomRole {
				return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
			}
			return errors.Wrap(err, "Custom role get error")
		}
	}

	if err := m.user.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	_, err := m.user.Get(ent.ID)
	if err != model.ErrNoSuchUser {
		m.user.AbortTx()
		if err == nil {
			return model.ErrUserAlreadyExists
		}
		return errors.Wrap(err, "Failed to get user info")
	}

	// Check duplicate user by name
	users, err := m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
	if err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to get user info by name")
	}
	if len(users) > 0 {
		m.user.AbortTx()
		return model.ErrUserAlreadyExists
	}

	if err := m.user.Add(ent); err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to add user")
	}
	m.user.CommitTx()
	return nil
}

// UserDelete ...
func (m *Manager) UserDelete(userID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	if err := m.user.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.user.Delete(userID); err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to delete user")
	}
	m.user.CommitTx()
	return nil
}

// UserGetList ...
func (m *Manager) UserGetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, error) {
	// TODO(validate projectName)
	return m.user.GetList(projectName, filter)
}

// UserGet ...
func (m *Manager) UserGet(userID string) (*model.UserInfo, error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.user.Get(userID)
}

// UserUpdate ...
func (m *Manager) UserUpdate(ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// Validate Role
	for _, r := range ent.SystemRoles {
		if !role.GetInst().IsValid(r) {
			return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
		}
	}
	for _, r := range ent.CustomRoles {
		if _, err := m.customRole.Get(r); err != nil {
			if errors.Cause(err) == model.ErrNoSuchCustomRole {
				return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
			}
			return errors.Wrap(err, "Custom role get error")
		}
	}

	if err := m.user.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}
	if err := m.user.Update(ent); err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to update user")
	}
	m.user.CommitTx()
	return nil
}

// UserAddRole ...
func (m *Manager) UserAddRole(userID string, roleType model.RoleType, roleID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	// Validate RoleID
	if roleType == model.RoleSystem {
		if !role.GetInst().IsValid(roleID) {
			return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
		}
	} else if roleType == model.RoleCustom {
		if _, err := m.customRole.Get(roleID); err != nil {
			if errors.Cause(err) == model.ErrNoSuchCustomRole {
				return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
			}
			return errors.Wrap(err, "Custom role get error")
		}
	}

	if err := m.user.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.user.AddRole(userID, roleType, roleID); err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to add role to user")
	}
	m.user.CommitTx()
	return nil
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(userID string, roleID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	// TODO(validate roleID)
	if err := m.user.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}
	if err := m.user.DeleteRole(userID, roleID); err != nil {
		m.user.AbortTx()
		return errors.Wrap(err, "Failed to delete role from user")
	}
	m.user.CommitTx()
	return nil
}

// UserSessionsDelete ...
func (m *Manager) UserSessionsDelete(userID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	if err := m.session.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.session.RevokeAll(userID); err != nil {
		m.session.AbortTx()
		return errors.Wrap(err, "Revoke session failed")
	}

	m.session.CommitTx()
	return nil
}

// LoginSessionAdd ...
func (m *Manager) LoginSessionAdd(info *model.LoginSessionInfo) error {
	// TODO(add validation)

	if err := m.loginSession.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}
	if err := m.loginSession.Add(info); err != nil {
		m.loginSession.AbortTx()
		return errors.Wrap(err, "Failed to add login session")
	}
	m.loginSession.CommitTx()
	return nil
}

// LoginSessionDelete ...
func (m *Manager) LoginSessionDelete(code string) error {
	if err := m.loginSession.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.loginSession.Delete(code); err != nil {
		m.loginSession.AbortTx()
		return errors.Wrap(err, "Failed to delete login session")
	}

	m.loginSession.CommitTx()
	return nil
}

// LoginSessionGet ...
func (m *Manager) LoginSessionGet(code string) (*model.LoginSessionInfo, error) {
	if err := m.loginSession.BeginTx(); err != nil {
		return nil, errors.Wrap(err, "BeginTx failed")
	}

	res, err := m.loginSession.Get(code)
	if err != nil {
		m.loginSession.AbortTx()
		return nil, errors.Wrap(err, "Failed to delete login session")
	}

	m.loginSession.CommitTx()
	return res, nil
}

// SessionAdd ...
func (m *Manager) SessionAdd(ent *model.Session) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if err := m.session.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if _, err := m.session.Get(ent.SessionID); err != model.ErrNoSuchSession {
		m.session.AbortTx()
		return model.ErrSessionAlreadyExists
	}

	if err := m.session.New(ent); err != nil {
		m.session.AbortTx()
		return errors.Wrap(err, "Failed to add session")
	}
	m.session.CommitTx()
	return nil
}

// SessionDelete ...
func (m *Manager) SessionDelete(sessionID string) error {
	if !model.ValidateSessionID(sessionID) {
		return errors.Wrap(model.ErrSessionValidateFailed, "invalid session id format")
	}

	if err := m.session.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.session.Revoke(sessionID); err != nil {
		m.session.AbortTx()
		return errors.Wrap(err, "Failed to revoke session")
	}
	m.session.CommitTx()
	return nil
}

// SessionGetList ...
func (m *Manager) SessionGetList(userID string) ([]string, error) {
	if !model.ValidateUserID(userID) {
		return []string{}, errors.Wrap(model.ErrSessionValidateFailed, "invalid user id format")
	}

	return m.session.GetList(userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if err := m.client.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	_, err := m.client.Get(ent.ID)
	if err != model.ErrNoSuchClient {
		m.client.AbortTx()
		if err == nil {
			return model.ErrClientAlreadyExists
		}
		return errors.Wrap(err, "Failed to get client info")
	}

	if err := m.client.Add(ent); err != nil {
		m.client.AbortTx()
		return errors.Wrap(err, "Failed to add client")
	}
	m.client.CommitTx()
	return nil
}

// ClientDelete ...
func (m *Manager) ClientDelete(clientID string) error {
	if !model.ValidateClientID(clientID) {
		return errors.Wrap(model.ErrClientValidateFailed, "invalid client id format")
	}

	if err := m.client.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.client.Delete(clientID); err != nil {
		m.client.AbortTx()
		return errors.Wrap(err, "Failed to delete client")
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
	if !model.ValidateClientID(clientID) {
		return nil, errors.Wrap(model.ErrClientValidateFailed, "invalid client id format")
	}

	return m.client.Get(clientID)
}

// ClientUpdate ...
func (m *Manager) ClientUpdate(ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if err := m.client.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.client.Update(ent); err != nil {
		m.client.AbortTx()
		return errors.Wrap(err, "Failed to update client")
	}
	m.client.CommitTx()
	return nil
}

// AuthCodeAdd ...
func (m *Manager) AuthCodeAdd(ent *model.AuthCode) error {
	// TODO(validate ent, identify by clientID and redirectURL)
	if err := m.authCode.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.authCode.New(ent); err != nil {
		m.authCode.AbortTx()
		return errors.Wrap(err, "Failed to add auth code")
	}
	m.authCode.CommitTx()
	return nil
}

// AuthCodeDelete ...
func (m *Manager) AuthCodeDelete(codeID string) error {
	// TODO(validate codeID)
	if err := m.authCode.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.authCode.Delete(codeID); err != nil {
		m.authCode.AbortTx()
		return errors.Wrap(err, "Failed to delete auth code")
	}
	m.authCode.CommitTx()
	return nil
}

// AuthCodeGet ...
func (m *Manager) AuthCodeGet(codeID string) (*model.AuthCode, error) {
	// TODO(validate codeID)
	return m.authCode.Get(codeID)
}

// CustomRoleAdd ...
func (m *Manager) CustomRoleAdd(ent *model.CustomRole) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// TODO(validate name uniquness in project)

	if err := m.customRole.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	_, err := m.customRole.Get(ent.ID)
	if err != model.ErrNoSuchCustomRole {
		m.customRole.AbortTx()
		if err == nil {
			return model.ErrCustomRoleAlreadyExists
		}
		return errors.Wrap(err, "Failed to get customRole info")
	}

	if err := m.customRole.Add(ent); err != nil {
		m.customRole.AbortTx()
		return errors.Wrap(err, "Failed to add customRole")
	}
	m.customRole.CommitTx()
	return nil
}

// CustomRoleDelete ...
func (m *Manager) CustomRoleDelete(customRoleID string) error {
	// TODO(validate customRoleID)

	if err := m.customRole.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.customRole.Delete(customRoleID); err != nil {
		m.customRole.AbortTx()
		return errors.Wrap(err, "Failed to delete customRole")
	}
	m.customRole.CommitTx()
	return nil
}

// CustomRoleGetList ...
func (m *Manager) CustomRoleGetList(projectName string) ([]string, error) {
	// TODO(validate projectName)
	return m.customRole.GetList(projectName)
}

// CustomRoleGet ...
func (m *Manager) CustomRoleGet(customRoleID string) (*model.CustomRole, error) {
	// TODO(validate customRoleID)
	return m.customRole.Get(customRoleID)
}

// CustomRoleUpdate ...
func (m *Manager) CustomRoleUpdate(ent *model.CustomRole) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// TODO(validate name uniquness in project)

	if err := m.customRole.BeginTx(); err != nil {
		return errors.Wrap(err, "BeginTx failed")
	}

	if err := m.customRole.Update(ent); err != nil {
		m.customRole.AbortTx()
		return errors.Wrap(err, "Failed to update client")
	}
	m.customRole.CommitTx()
	return nil
}
