package db

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/memory"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/pwpol"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

// Manager ...
type Manager struct {
	project         model.ProjectInfoHandler
	user            model.UserInfoHandler
	session         model.SessionHandler
	client          model.ClientInfoHandler
	customRole      model.CustomRoleHandler
	authCodeSession model.AuthCodeSessionHandler
	transaction     model.TransactionManager
	ping            model.PingHandler
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
		inst = &Manager{
			project:         memory.NewProjectHandler(),
			user:            memory.NewUserHandler(),
			session:         memory.NewSessionHandler(),
			client:          memory.NewClientHandler(),
			customRole:      memory.NewCustomRoleHandler(),
			authCodeSession: memory.NewAuthCodeSessionHandler(),
			transaction:     memory.NewTransactionManager(),
			ping:            memory.NewPingHandler(),
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
		clientHandler, err := mongo.NewClientHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create client handler")
		}
		userHandler, err := mongo.NewUserHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create user handler")
		}
		sessionHandler, err := mongo.NewSessionHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create session handler")
		}
		customRoleHandler, err := mongo.NewCustomRoleHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create custom role handler")
		}
		authCodeSessionHandler, err := mongo.NewAuthCodeSessionHandler(dbClient)
		if err != nil {
			return errors.Wrap(err, "Failed to create login session handler")
		}

		inst = &Manager{
			project:         prjHandler,
			user:            userHandler,
			session:         sessionHandler,
			client:          clientHandler,
			customRole:      customRoleHandler,
			authCodeSession: authCodeSessionHandler,
			transaction:     mongo.NewTransactionManager(dbClient),
			ping:            mongo.NewPingHandler(dbClient),
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

// Ping ...
func (m *Manager) Ping() error {
	return m.ping.Ping()
}

// ProjectAdd ...
func (m *Manager) ProjectAdd(ent *model.ProjectInfo) error {
	if err := ent.Validate(); err != nil {
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

	return m.transaction.Transaction(func() error {
		if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
			return model.ErrProjectAlreadyExists
		}

		if err := m.project.Add(ent); err != nil {
			return errors.Wrap(err, "Failed to add project")
		}

		callbacks := []string{
			"http://localhost:3000/callback", // TODO(for debug)
		}
		if os.Getenv("HEKATE_PORTAL_ADDR") != "" {
			addr := os.Getenv("HEKATE_PORTAL_ADDR") + "/callback"
			callbacks = append(callbacks, addr)
		}
		// add client for portal login
		clientEnt := &model.ClientInfo{
			ID:                  "portal",
			ProjectName:         ent.Name,
			AccessType:          "public",
			CreatedAt:           ent.CreatedAt,
			AllowedCallbackURLs: callbacks,
		}
		if err := m.client.Add(ent.Name, clientEnt); err != nil {
			return errors.Wrap(err, "Failed to add client for portal login")
		}

		return nil
	})
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) error {
	if name == "" {
		return errors.Wrap(model.ErrProjectValidateFailed, "name of entry is empty")
	}

	return m.transaction.Transaction(func() error {
		prj, err := m.project.Get(name)
		if err != nil {
			return errors.Wrap(err, "Failed to get delete project info")
		}

		if !prj.PermitDelete {
			return errors.Wrap(model.ErrDeleteBlockedProject, "the project can not delete")
		}

		if err := m.authCodeSession.DeleteAllInProject(name); err != nil {
			return errors.Wrap(err, "Failed to delete login session data")
		}

		if err := m.session.DeleteAllInProject(name); err != nil {
			return errors.Wrap(err, "Failed to delete session data")
		}

		if err := m.customRole.DeleteAll(name); err != nil {
			return errors.Wrap(err, "Failed to delete custom role data")
		}

		if err := m.client.DeleteAll(name); err != nil {
			return errors.Wrap(err, "Failed to delete client data")
		}

		if err := m.user.DeleteAll(name); err != nil {
			return errors.Wrap(err, "Failed to delete user data")
		}

		if err := m.project.Delete(name); err != nil {
			return errors.Wrap(err, "Failed to delete project")
		}

		return nil
	})
}

// ProjectGetList ...
func (m *Manager) ProjectGetList() ([]*model.ProjectInfo, error) {
	return m.project.GetList()
}

// ProjectGet ...
func (m *Manager) ProjectGet(name string) (*model.ProjectInfo, error) {
	if name == "" {
		return nil, errors.Wrap(model.ErrProjectValidateFailed, "name of entry is empty")
	}

	return m.project.Get(name)
}

// ProjectUpdate ...
func (m *Manager) ProjectUpdate(ent *model.ProjectInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate")
	}

	return m.transaction.Transaction(func() error {
		if err := m.project.Update(ent); err != nil {
			return errors.Wrap(err, "Failed to update project")
		}
		return nil
	})
}

// UserAdd ...
func (m *Manager) UserAdd(projectName string, ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// Validate Roles
	for _, r := range ent.SystemRoles {
		res, typ, ok := role.GetInst().Parse(r)
		if !ok {
			return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
		}

		// Require read permission if append write permission
		if *typ == role.TypeWrite {
			if ok := role.Authorize(ent.SystemRoles, *res, role.TypeRead); !ok {
				return errors.Wrap(model.ErrUserValidateFailed, "Do not have read permission")
			}
		}
	}

	return m.transaction.Transaction(func() error {
		for _, r := range ent.CustomRoles {
			if _, err := m.customRole.Get(projectName, r); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		_, err := m.user.Get(projectName, ent.ID)
		if err != model.ErrNoSuchUser {
			if err == nil {
				return model.ErrUserAlreadyExists
			}
			return errors.Wrap(err, "Failed to get user info")
		}

		// Check duplicate user by name
		users, err := m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
		if err != nil {
			return errors.Wrap(err, "Failed to get user info by name")
		}
		if len(users) > 0 {
			return model.ErrUserAlreadyExists
		}

		if err := m.user.Add(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to add user")
		}
		return nil
	})
}

// UserDelete ...
func (m *Manager) UserDelete(projectName string, userID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.authCodeSession.DeleteAllInUser(projectName, userID); err != nil {
			return errors.Wrap(err, "Delete authoriation code failed")
		}

		if err := m.session.DeleteAll(projectName, userID); err != nil {
			return errors.Wrap(err, "Delete user session failed")
		}

		if err := m.user.Delete(projectName, userID); err != nil {
			return errors.Wrap(err, "Failed to delete user")
		}
		return nil
	})
}

// UserGetList ...
func (m *Manager) UserGetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, error) {
	if !model.ValidateProjectName(projectName) {
		return nil, errors.Wrap(model.ErrUserValidateFailed, "invalid project name format")
	}

	return m.user.GetList(projectName, filter)
}

// UserGet ...
func (m *Manager) UserGet(projectName string, userID string) (*model.UserInfo, error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.user.Get(projectName, userID)
}

// UserUpdate ...
func (m *Manager) UserUpdate(projectName string, ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// Validate Role
	for _, r := range ent.SystemRoles {
		res, typ, ok := role.GetInst().Parse(r)
		if !ok {
			return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
		}

		// Require read permission if append write permission
		if *typ == role.TypeWrite {
			if ok := role.Authorize(ent.SystemRoles, *res, role.TypeRead); !ok {
				return errors.Wrap(model.ErrUserValidateFailed, "Do not have read permission")
			}
		}
	}

	return m.transaction.Transaction(func() error {
		for _, r := range ent.CustomRoles {
			if _, err := m.customRole.Get(projectName, r); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		// check duplicate user name
		users, err := m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
		if err != nil {
			return errors.Wrap(err, "Failed to get user for checking name duplication")
		}
		if len(users) >= 2 || (len(users) == 1 && users[0].ID != ent.ID) {
			return errors.Wrap(model.ErrUserAlreadyExists, "new user name is already used")
		}

		if err := m.user.Update(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to update user")
		}
		return nil
	})
}

// UserAddRole ...
func (m *Manager) UserAddRole(projectName string, userID string, roleType model.RoleType, roleID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		// Validate RoleID
		if roleType == model.RoleSystem {
			res, typ, ok := role.GetInst().Parse(roleID)
			if !ok {
				return errors.Wrap(model.ErrUserValidateFailed, "Invalid system role")
			}

			usr, err := m.user.Get(projectName, userID)

			if *res == role.ResCluster && usr.ProjectName != "master" {
				return errors.Wrap(model.ErrUserValidateFailed, "Resource cluster can add to master project user")
			}

			// check user already has read permission if roleID type is write
			if *typ == role.TypeWrite {
				if err != nil {
					return errors.Wrap(err, "Failed to get user system roles")
				}

				// If user do not have read permission to the target resource, return error
				if !role.Authorize(usr.SystemRoles, *res, role.TypeRead) {
					return errors.Wrap(model.ErrUserValidateFailed, "Do not have read permission")
				}
			}
		} else if roleType == model.RoleCustom {
			if _, err := m.customRole.Get(projectName, roleID); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		if err := m.user.AddRole(projectName, userID, roleType, roleID); err != nil {
			return errors.Wrap(err, "Failed to add role to user")
		}
		return nil
	})
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(projectName string, userID string, roleID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		usr, err := m.user.Get(projectName, userID)
		if err != nil {
			return errors.Wrap(err, "Failed to get user system roles")
		}

		res, typ, ok := role.GetInst().Parse(roleID)
		if ok {
			// roleID is system role, so check write permission
			if *typ == role.TypeRead {
				if role.Authorize(usr.SystemRoles, *res, role.TypeWrite) {
					return errors.Wrap(model.ErrUserValidateFailed, "Remove write permission at first")
				}
			}
		}
		// If not ok, roleID maybe Custom Role

		if err := m.user.DeleteRole(projectName, userID, roleID); err != nil {
			return errors.Wrap(err, "Failed to delete role from user")
		}
		return nil
	})
}

// UserChangePassword ...
func (m *Manager) UserChangePassword(projectName string, userID string, password string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		usr, err := m.user.Get(projectName, userID)
		if err != nil {
			return errors.Wrap(err, "Failed to get user of change password")
		}

		prj, err := m.project.Get(projectName)
		if err != nil {
			return errors.Wrap(err, "Failed to get project associated with the user")
		}

		if err := pwpol.CheckPassword(usr.Name, password, prj.PasswordPolicy); err != nil {
			return errors.Wrap(err, "Failed to check password")
		}

		usr.PasswordHash = util.CreateHash(password)

		if err := m.user.Update(projectName, usr); err != nil {
			return errors.Wrap(err, "Failed to update user password")
		}

		return nil
	})
}

// AuthCodeSessionAdd ...
func (m *Manager) AuthCodeSessionAdd(projectName string, ent *model.AuthCodeSession) error {
	// create session is in internal only, so validation is not required
	return m.transaction.Transaction(func() error {
		if err := m.authCodeSession.Add(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to add auth code session")
		}
		return nil
	})
}

// AuthCodeSessionUpdate ...
func (m *Manager) AuthCodeSessionUpdate(projectName string, ent *model.AuthCodeSession) error {
	// update session is in internal only, so validation is not required
	return m.transaction.Transaction(func() error {
		if err := m.authCodeSession.Update(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to update auth code session")
		}
		return nil
	})
}

// AuthCodeSessionDelete ...
func (m *Manager) AuthCodeSessionDelete(projectName string, sessionID string) error {
	return m.transaction.Transaction(func() error {
		if err := m.authCodeSession.Delete(projectName, sessionID); err != nil {
			return errors.Wrap(err, "Failed to delete auth code session")
		}

		return nil
	})
}

// AuthCodeSessionGet ...
func (m *Manager) AuthCodeSessionGet(projectName string, sessionID string) (*model.AuthCodeSession, error) {
	return m.authCodeSession.Get(projectName, sessionID)
}

// AuthCodeSessionGetByCode ...
func (m *Manager) AuthCodeSessionGetByCode(projectName string, code string) (*model.AuthCodeSession, error) {
	return m.authCodeSession.GetByCode(projectName, code)
}

// SessionAdd ...
func (m *Manager) SessionAdd(projectName string, ent *model.Session) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		if _, err := m.session.Get(projectName, ent.SessionID); err != model.ErrNoSuchSession {
			return model.ErrSessionAlreadyExists
		}

		if err := m.session.Add(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to add session")
		}
		return nil
	})
}

// SessionGet ..
func (m *Manager) SessionGet(projectName string, sessionID string) (*model.Session, error) {
	// TODO(add validation)
	return m.session.Get(projectName, sessionID)
}

// SessionDelete ...
func (m *Manager) SessionDelete(projectName string, sessionID string) error {
	if !model.ValidateSessionID(sessionID) {
		return errors.Wrap(model.ErrSessionValidateFailed, "invalid session id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.session.Delete(projectName, sessionID); err != nil {
			return errors.Wrap(err, "Failed to revoke session")
		}
		return nil
	})
}

// SessionGetList ...
func (m *Manager) SessionGetList(projectName string, userID string) ([]*model.Session, error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Wrap(model.ErrSessionValidateFailed, "invalid user id format")
	}

	return m.session.GetList(projectName, userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(projectName string, ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		_, err := m.client.Get(ent.ProjectName, ent.ID)
		if err != model.ErrNoSuchClient {
			if err == nil {
				return model.ErrClientAlreadyExists
			}
			return errors.Wrap(err, "Failed to get client info")
		}

		if err := m.client.Add(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to add client")
		}
		return nil
	})
}

// ClientDelete ...
func (m *Manager) ClientDelete(projectName, clientID string) error {
	if !model.ValidateProjectName(projectName) {
		return errors.Wrap(model.ErrClientValidateFailed, "Invalid project name format")
	}
	if !model.ValidateClientID(clientID) {
		return errors.Wrap(model.ErrClientValidateFailed, "invalid client id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.authCodeSession.DeleteAllInClient(projectName, clientID); err != nil {
			return errors.Wrap(err, "Failed to delete login session of the client")
		}

		if err := m.client.Delete(projectName, clientID); err != nil {
			return errors.Wrap(err, "Failed to delete client")
		}
		return nil
	})
}

// ClientGetList ...
func (m *Manager) ClientGetList(projectName string) ([]*model.ClientInfo, error) {
	if !model.ValidateProjectName(projectName) {
		return nil, errors.Wrap(model.ErrClientValidateFailed, "Invalid project name format")
	}

	return m.client.GetList(projectName)
}

// ClientGet ...
func (m *Manager) ClientGet(projectName, clientID string) (*model.ClientInfo, error) {
	if !model.ValidateProjectName(projectName) {
		return nil, errors.Wrap(model.ErrClientValidateFailed, "Invalid project name format")
	}
	if !model.ValidateClientID(clientID) {
		return nil, errors.Wrap(model.ErrClientValidateFailed, "invalid client id format")
	}

	return m.client.Get(projectName, clientID)
}

// ClientUpdate ...
func (m *Manager) ClientUpdate(projectName string, ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		if err := m.client.Update(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to update client")
		}
		return nil
	})
}

// CustomRoleAdd ...
func (m *Manager) CustomRoleAdd(projectName string, ent *model.CustomRole) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	// TODO(validate name uniquness in project)

	return m.transaction.Transaction(func() error {
		_, err := m.customRole.Get(projectName, ent.ID)
		if err != model.ErrNoSuchCustomRole {
			if err == nil {
				return model.ErrCustomRoleAlreadyExists
			}
			return errors.Wrap(err, "Failed to get customRole info")
		}

		if err := m.customRole.Add(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to add customRole")
		}
		return nil
	})
}

// CustomRoleDelete ...
func (m *Manager) CustomRoleDelete(projectName string, customRoleID string) error {
	// TODO(validate customRoleID)

	return m.transaction.Transaction(func() error {
		if err := m.user.DeleteAllCustomRole(projectName, customRoleID); err != nil {
			return errors.Wrap(err, "Failed to delete custom role from user")
		}

		if err := m.customRole.Delete(projectName, customRoleID); err != nil {
			return errors.Wrap(err, "Failed to delete customRole")
		}
		return nil
	})
}

// CustomRoleGetList ...
func (m *Manager) CustomRoleGetList(projectName string, filter *model.CustomRoleFilter) ([]*model.CustomRole, error) {
	if !model.ValidateProjectName(projectName) {
		return nil, errors.Wrap(model.ErrCustomRoleValidateFailed, "invalid project name format")
	}
	return m.customRole.GetList(projectName, filter)
}

// CustomRoleGet ...
func (m *Manager) CustomRoleGet(projectName string, customRoleID string) (*model.CustomRole, error) {
	// TODO(validate customRoleID)
	return m.customRole.Get(projectName, customRoleID)
}

// CustomRoleUpdate ...
func (m *Manager) CustomRoleUpdate(projectName string, ent *model.CustomRole) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		r, err := m.customRole.GetList(ent.ProjectName, &model.CustomRoleFilter{Name: ent.Name})
		if err != nil {
			return errors.Wrap(err, "Failed to get role list")
		}

		// check name uniquness in project
		if len(r) > 0 && r[0].ID != ent.ID {
			return model.ErrCustomRoleAlreadyExists
		}

		if err := m.customRole.Update(projectName, ent); err != nil {
			return errors.Wrap(err, "Failed to update client")
		}
		return nil
	})
}
