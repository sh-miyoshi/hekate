package db

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/memory"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/util"
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
	transaction  model.TransactionManager
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
			project:      memory.NewProjectHandler(),
			user:         memory.NewUserHandler(),
			session:      memory.NewSessionHandler(),
			client:       memory.NewClientHandler(),
			authCode:     memory.NewAuthCodeHandler(),
			customRole:   memory.NewCustomRoleHandler(),
			loginSession: memory.NewLoginSessionHandler(),
			transaction:  memory.NewTransactionManager(),
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
			transaction:  mongo.NewTransactionManager(dbClient),
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

	return m.transaction.Transaction(func() error {
		if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
			return model.ErrProjectAlreadyExists
		}

		if err := m.project.Add(ent); err != nil {
			return errors.Wrap(err, "Failed to add project")
		}

		return nil
	})
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) error {
	if name == "" {
		return errors.New("name of entry is empty")
	}

	return m.transaction.Transaction(func() error {
		prj, err := m.project.Get(name)
		if err != nil {
			return errors.Wrap(err, "Failed to get delete project info")
		}

		if !prj.PermitDelete {
			return errors.Wrap(model.ErrDeleteBlockedProject, "the project can not delete")
		}

		// TODO(delete loginsession, session)

		if err := m.authCode.DeleteAll(name); err != nil {
			return errors.Wrap(err, "Failed to delete oidc code data")
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

	return m.transaction.Transaction(func() error {
		if err := m.project.Update(ent); err != nil {
			return errors.Wrap(err, "Failed to update project")
		}
		return nil
	})
}

// UserAdd ...
func (m *Manager) UserAdd(ent *model.UserInfo) error {
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
			if _, err := m.customRole.Get(r); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		_, err := m.user.Get(ent.ID)
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

		if err := m.user.Add(ent); err != nil {
			return errors.Wrap(err, "Failed to add user")
		}
		return nil
	})
}

// UserDelete ...
func (m *Manager) UserDelete(userID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.authCode.DeleteAll(userID); err != nil {
			return errors.Wrap(err, "Delete authoriation code failed")
		}

		if err := m.session.RevokeAll(userID); err != nil {
			return errors.Wrap(err, "Delete user session failed")
		}

		if err := m.user.Delete(userID); err != nil {
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
			if _, err := m.customRole.Get(r); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		if err := m.user.Update(ent); err != nil {
			return errors.Wrap(err, "Failed to update user")
		}
		return nil
	})
}

// UserAddRole ...
func (m *Manager) UserAddRole(userID string, roleType model.RoleType, roleID string) error {
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

			usr, err := m.user.Get(userID)

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
			if _, err := m.customRole.Get(roleID); err != nil {
				if errors.Cause(err) == model.ErrNoSuchCustomRole {
					return errors.Wrap(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Wrap(err, "Custom role get error")
			}
		}

		if err := m.user.AddRole(userID, roleType, roleID); err != nil {
			return errors.Wrap(err, "Failed to add role to user")
		}
		return nil
	})
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(userID string, roleID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		usr, err := m.user.Get(userID)
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

		if err := m.user.DeleteRole(userID, roleID); err != nil {
			return errors.Wrap(err, "Failed to delete role from user")
		}
		return nil
	})
}

// UserSessionsDelete ...
func (m *Manager) UserSessionsDelete(userID string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.session.RevokeAll(userID); err != nil {
			return errors.Wrap(err, "Revoke session failed")
		}

		return nil
	})
}

// UserChangePassword ...
func (m *Manager) UserChangePassword(userID string, password string) error {
	if !model.ValidateUserID(userID) {
		return errors.Wrap(model.ErrUserValidateFailed, "invalid user id format")
	}

	// TODO(validate password)

	return m.transaction.Transaction(func() error {
		usr, err := m.user.Get(userID)
		if err != nil {
			return errors.Wrap(err, "Failed to get user of change password")
		}

		usr.PasswordHash = util.CreateHash(password)

		if err := m.user.Update(usr); err != nil {
			return errors.Wrap(err, "Failed to update user password")
		}

		return nil
	})
}

// LoginSessionAdd ...
func (m *Manager) LoginSessionAdd(info *model.LoginSessionInfo) error {
	// TODO(add validation)

	return m.transaction.Transaction(func() error {
		if err := m.loginSession.Add(info); err != nil {
			return errors.Wrap(err, "Failed to add login session")
		}
		return nil
	})
}

// LoginSessionDelete ...
func (m *Manager) LoginSessionDelete(code string) error {
	return m.transaction.Transaction(func() error {
		if err := m.loginSession.Delete(code); err != nil {
			return errors.Wrap(err, "Failed to delete login session")
		}

		return nil
	})
}

// LoginSessionGet ...
func (m *Manager) LoginSessionGet(code string) (*model.LoginSessionInfo, error) {
	//TODO(add code validation)
	return m.loginSession.Get(code)
}

// SessionAdd ...
func (m *Manager) SessionAdd(ent *model.Session) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		if _, err := m.session.Get(ent.SessionID); err != model.ErrNoSuchSession {
			return model.ErrSessionAlreadyExists
		}

		if err := m.session.New(ent); err != nil {
			return errors.Wrap(err, "Failed to add session")
		}
		return nil
	})
}

// SessionDelete ...
func (m *Manager) SessionDelete(sessionID string) error {
	if !model.ValidateSessionID(sessionID) {
		return errors.Wrap(model.ErrSessionValidateFailed, "invalid session id format")
	}

	return m.transaction.Transaction(func() error {
		if err := m.session.Revoke(sessionID); err != nil {
			return errors.Wrap(err, "Failed to revoke session")
		}
		return nil
	})
}

// SessionGetList ...
func (m *Manager) SessionGetList(userID string) ([]*model.Session, error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Wrap(model.ErrSessionValidateFailed, "invalid user id format")
	}

	return m.session.GetList(userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(ent *model.ClientInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() error {
		_, err := m.client.Get(ent.ID)
		if err != model.ErrNoSuchClient {
			if err == nil {
				return model.ErrClientAlreadyExists
			}
			return errors.Wrap(err, "Failed to get client info")
		}

		if err := m.client.Add(ent); err != nil {
			return errors.Wrap(err, "Failed to add client")
		}
		return nil
	})
}

// ClientDelete ...
func (m *Manager) ClientDelete(clientID string) error {
	if !model.ValidateClientID(clientID) {
		return errors.Wrap(model.ErrClientValidateFailed, "invalid client id format")
	}

	return m.transaction.Transaction(func() error {

		// TODO(delete loginsession, oidc_code, session)

		if err := m.client.Delete(clientID); err != nil {
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

	return m.transaction.Transaction(func() error {
		if err := m.client.Update(ent); err != nil {
			return errors.Wrap(err, "Failed to update client")
		}
		return nil
	})
}

// AuthCodeAdd ...
func (m *Manager) AuthCodeAdd(ent *model.AuthCode) error {
	// TODO(validate ent, identify by clientID and redirectURL)
	return m.transaction.Transaction(func() error {
		if err := m.authCode.New(ent); err != nil {
			return errors.Wrap(err, "Failed to add auth code")
		}
		return nil
	})
}

// AuthCodeDelete ...
func (m *Manager) AuthCodeDelete(codeID string) error {
	// TODO(validate codeID)
	return m.transaction.Transaction(func() error {
		if err := m.authCode.Delete(codeID); err != nil {
			return errors.Wrap(err, "Failed to delete auth code")
		}
		return nil
	})
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

	return m.transaction.Transaction(func() error {
		_, err := m.customRole.Get(ent.ID)
		if err != model.ErrNoSuchCustomRole {
			if err == nil {
				return model.ErrCustomRoleAlreadyExists
			}
			return errors.Wrap(err, "Failed to get customRole info")
		}

		if err := m.customRole.Add(ent); err != nil {
			return errors.Wrap(err, "Failed to add customRole")
		}
		return nil
	})
}

// CustomRoleDelete ...
func (m *Manager) CustomRoleDelete(customRoleID string) error {
	// TODO(validate customRoleID)

	return m.transaction.Transaction(func() error {
		// TODO(delete role from user)

		if err := m.customRole.Delete(customRoleID); err != nil {
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

	return m.transaction.Transaction(func() error {
		if err := m.customRole.Update(ent); err != nil {
			return errors.Wrap(err, "Failed to update client")
		}
		return nil
	})
}
