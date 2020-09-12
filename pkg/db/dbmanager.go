package db

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/sh-miyoshi/hekate/pkg/db/memory"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/pwpol"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

// Manager ...
type Manager struct {
	project      model.ProjectInfoHandler
	user         model.UserInfoHandler
	session      model.SessionHandler
	client       model.ClientInfoHandler
	customRole   model.CustomRoleHandler
	loginSession model.LoginSessionHandler
	transaction  model.TransactionManager
	ping         model.PingHandler

	portalAddr string
}

var inst *Manager

// InitDBManager ...
func InitDBManager(dbType string, connStr string) *errors.Error {
	if inst != nil {
		return errors.New("Internal server error", "DBManager is already initialized")
	}

	switch dbType {
	case "memory":
		logger.Info("Initialize with local memory DB")
		inst = &Manager{
			project:      memory.NewProjectHandler(),
			user:         memory.NewUserHandler(),
			session:      memory.NewSessionHandler(),
			client:       memory.NewClientHandler(),
			customRole:   memory.NewCustomRoleHandler(),
			loginSession: memory.NewLoginSessionHandler(),
			transaction:  memory.NewTransactionManager(),
			ping:         memory.NewPingHandler(),
		}
	case "mongo":
		logger.Info("Initialize with mongo DB")
		dbClient, err := mongo.NewClient(connStr)
		if err != nil {
			return errors.Append(err, "Failed to create db client")
		}

		prjHandler, err := mongo.NewProjectHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create project handler")
		}
		clientHandler, err := mongo.NewClientHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create client handler")
		}
		userHandler, err := mongo.NewUserHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create user handler")
		}
		sessionHandler, err := mongo.NewSessionHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create session handler")
		}
		customRoleHandler, err := mongo.NewCustomRoleHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create custom role handler")
		}
		loginSessionHandler, err := mongo.NewLoginSessionHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create login session handler")
		}

		inst = &Manager{
			project:      prjHandler,
			user:         userHandler,
			session:      sessionHandler,
			client:       clientHandler,
			customRole:   customRoleHandler,
			loginSession: loginSessionHandler,
			transaction:  mongo.NewTransactionManager(dbClient),
			ping:         mongo.NewPingHandler(dbClient),
		}
	default:
		return errors.New("Internal server error", "Database Type %s is not implemented yet", dbType)
	}

	if os.Getenv("HEKATE_PORTAL_ADDR") != "" {
		inst.portalAddr = os.Getenv("HEKATE_PORTAL_ADDR") + "/callback"
		if !govalidator.IsRequestURL(inst.portalAddr) {
			return errors.Append(model.ErrProjectValidateFailed, "Invalid portal callback URL %s is specified.", inst.portalAddr)
		}
		logger.Info("Set portal callback URL: %s", inst.portalAddr)
	}

	return nil
}

// GetInst returns an instance of DB Manager
func GetInst() *Manager {
	return inst
}

// Ping ...
func (m *Manager) Ping() *errors.Error {
	return m.ping.Ping()
}

// ProjectAdd ...
func (m *Manager) ProjectAdd(ent *model.ProjectInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Validate failed")
	}

	switch ent.TokenConfig.SigningAlgorithm {
	case "RS256":
		key, err := rsa.GenerateKey(rand.Reader, 2048) // fixed key length is ok?
		if err != nil {
			return errors.New("RSA key generate failed", "Failed to generate RSA private key: %v", err)
		}
		ent.TokenConfig.SignSecretKey = x509.MarshalPKCS1PrivateKey(key)
		ent.TokenConfig.SignPublicKey = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	}

	return m.transaction.Transaction(func() *errors.Error {
		if _, err := m.project.Get(ent.Name); err != model.ErrNoSuchProject {
			return model.ErrProjectAlreadyExists
		}

		if err := m.project.Add(ent); err != nil {
			return errors.Append(err, "Failed to add project")
		}

		callbacks := []string{}
		if m.portalAddr != "" {
			callbacks = append(callbacks, m.portalAddr)
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
			return errors.Append(err, "Failed to add client for portal login")
		}

		return nil
	})
}

// ProjectDelete ...
func (m *Manager) ProjectDelete(name string) *errors.Error {
	if !model.ValidateProjectName(name) {
		return errors.Append(model.ErrProjectValidateFailed, "Invalid project name format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		prj, err := m.project.Get(name)
		if err != nil {
			return errors.Append(err, "Failed to get delete project info")
		}

		if !prj.PermitDelete {
			return errors.Append(model.ErrDeleteBlockedProject, "the project can not delete")
		}

		if err := m.loginSession.DeleteAllInProject(name); err != nil {
			return errors.Append(err, "Failed to delete login session data")
		}

		if err := m.session.DeleteAllInProject(name); err != nil {
			return errors.Append(err, "Failed to delete session data")
		}

		if err := m.customRole.DeleteAll(name); err != nil {
			return errors.Append(err, "Failed to delete custom role data")
		}

		if err := m.client.DeleteAll(name); err != nil {
			return errors.Append(err, "Failed to delete client data")
		}

		if err := m.user.DeleteAll(name); err != nil {
			return errors.Append(err, "Failed to delete user data")
		}

		if err := m.project.Delete(name); err != nil {
			return errors.Append(err, "Failed to delete project")
		}

		return nil
	})
}

// ProjectGetList ...
func (m *Manager) ProjectGetList() ([]*model.ProjectInfo, *errors.Error) {
	return m.project.GetList()
}

// ProjectGet ...
func (m *Manager) ProjectGet(name string) (*model.ProjectInfo, *errors.Error) {
	if !model.ValidateProjectName(name) {
		return nil, errors.Append(model.ErrProjectValidateFailed, "Invalid project name format")
	}

	return m.project.Get(name)
}

// ProjectUpdate ...
func (m *Manager) ProjectUpdate(ent *model.ProjectInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.project.Update(ent); err != nil {
			return errors.Append(err, "Failed to update project")
		}
		return nil
	})
}

// UserAdd ...
func (m *Manager) UserAdd(projectName string, ent *model.UserInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	// Validate Roles
	for _, r := range ent.SystemRoles {
		res, typ, ok := role.GetInst().Parse(r)
		if !ok {
			return errors.Append(model.ErrUserValidateFailed, "Invalid system role")
		}

		// Require read permission if append write permission
		if *typ == role.TypeWrite {
			if ok := role.Authorize(ent.SystemRoles, *res, role.TypeRead); !ok {
				return errors.Append(model.ErrUserValidateFailed, "Do not have read permission")
			}
		}
	}

	return m.transaction.Transaction(func() *errors.Error {
		for _, r := range ent.CustomRoles {
			if _, err := m.customRole.Get(projectName, r); err != nil {
				if errors.Contains(err, model.ErrNoSuchCustomRole) {
					return errors.Append(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Append(err, "Custom role get error")
			}
		}

		_, err := m.user.Get(projectName, ent.ID)
		if err != model.ErrNoSuchUser {
			if err == nil {
				return model.ErrUserAlreadyExists
			}
			return errors.Append(err, "Failed to get user info")
		}

		// Check duplicate user by name
		users, err := m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
		if err != nil {
			return errors.Append(err, "Failed to get user info by name")
		}
		if len(users) > 0 {
			return model.ErrUserAlreadyExists
		}

		if err := m.user.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add user")
		}
		return nil
	})
}

// UserDelete ...
func (m *Manager) UserDelete(projectName string, userID string) *errors.Error {
	if !model.ValidateUserID(userID) {
		return errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.DeleteAllInUser(projectName, userID); err != nil {
			return errors.Append(err, "Delete authoriation code failed")
		}

		if err := m.session.DeleteAll(projectName, userID); err != nil {
			return errors.Append(err, "Delete user session failed")
		}

		if err := m.user.Delete(projectName, userID); err != nil {
			return errors.Append(err, "Failed to delete user")
		}
		return nil
	})
}

// UserGetList ...
func (m *Manager) UserGetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, *errors.Error) {
	return m.user.GetList(projectName, filter)
}

// UserGet ...
func (m *Manager) UserGet(projectName string, userID string) (*model.UserInfo, *errors.Error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.user.Get(projectName, userID)
}

// UserUpdate ...
func (m *Manager) UserUpdate(projectName string, ent *model.UserInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	// Validate Role
	for _, r := range ent.SystemRoles {
		res, typ, ok := role.GetInst().Parse(r)
		if !ok {
			return errors.Append(model.ErrUserValidateFailed, "Invalid system role")
		}

		// Require read permission if append write permission
		if *typ == role.TypeWrite {
			if ok := role.Authorize(ent.SystemRoles, *res, role.TypeRead); !ok {
				return errors.Append(model.ErrUserValidateFailed, "Do not have read permission")
			}
		}
	}

	return m.transaction.Transaction(func() *errors.Error {
		for _, r := range ent.CustomRoles {
			if _, err := m.customRole.Get(projectName, r); err != nil {
				if errors.Contains(err, model.ErrNoSuchCustomRole) {
					return errors.Append(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Append(err, "Custom role get error")
			}
		}

		// check duplicate user name
		users, err := m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
		if err != nil {
			return errors.Append(err, "Failed to get user for checking name duplication")
		}
		if len(users) >= 2 || (len(users) == 1 && users[0].ID != ent.ID) {
			return errors.Append(model.ErrUserAlreadyExists, "new user name is already used")
		}

		if err := m.user.Update(projectName, ent); err != nil {
			return errors.Append(err, "Failed to update user")
		}
		return nil
	})
}

// UserAddRole ...
func (m *Manager) UserAddRole(projectName string, userID string, roleType model.RoleType, roleID string) *errors.Error {
	if !model.ValidateUserID(userID) {
		return errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		// Validate RoleID
		if roleType == model.RoleSystem {
			res, typ, ok := role.GetInst().Parse(roleID)
			if !ok {
				return errors.Append(model.ErrUserValidateFailed, "Invalid system role")
			}

			usr, err := m.user.Get(projectName, userID)

			if *res == role.ResCluster && usr.ProjectName != "master" {
				return errors.Append(model.ErrUserValidateFailed, "Resource cluster can add to master project user")
			}

			// check user already has read permission if roleID type is write
			if *typ == role.TypeWrite {
				if err != nil {
					return errors.Append(err, "Failed to get user system roles")
				}

				// If user do not have read permission to the target resource, return error
				if !role.Authorize(usr.SystemRoles, *res, role.TypeRead) {
					return errors.Append(model.ErrUserValidateFailed, "Do not have read permission")
				}
			}
		} else if roleType == model.RoleCustom {
			if _, err := m.customRole.Get(projectName, roleID); err != nil {
				if errors.Contains(err, model.ErrNoSuchCustomRole) {
					return errors.Append(model.ErrUserValidateFailed, "Invalid custom role")
				}
				return errors.Append(err, "Custom role get error")
			}
		}

		if err := m.user.AddRole(projectName, userID, roleType, roleID); err != nil {
			return errors.Append(err, "Failed to add role to user")
		}
		return nil
	})
}

// UserDeleteRole ...
func (m *Manager) UserDeleteRole(projectName string, userID string, roleID string) *errors.Error {
	if !model.ValidateUserID(userID) {
		return errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		usr, err := m.user.Get(projectName, userID)
		if err != nil {
			return errors.Append(err, "Failed to get user system roles")
		}

		res, typ, ok := role.GetInst().Parse(roleID)
		if ok {
			// roleID is system role, so check write permission
			if *typ == role.TypeRead {
				if role.Authorize(usr.SystemRoles, *res, role.TypeWrite) {
					return errors.Append(model.ErrUserValidateFailed, "Remove write permission at first")
				}
			}
		}
		// If not ok, roleID maybe Custom Role

		if err := m.user.DeleteRole(projectName, userID, roleID); err != nil {
			return errors.Append(err, "Failed to delete role from user")
		}
		return nil
	})
}

// UserChangePassword ...
func (m *Manager) UserChangePassword(projectName string, userID string, password string) *errors.Error {
	if !model.ValidateUserID(userID) {
		return errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		usr, err := m.user.Get(projectName, userID)
		if err != nil {
			return errors.Append(err, "Failed to get user of change password")
		}

		prj, err := m.project.Get(projectName)
		if err != nil {
			return errors.Append(err, "Failed to get project associated with the user")
		}

		if err := pwpol.CheckPassword(usr.Name, password, prj.PasswordPolicy); err != nil {
			return errors.Append(err, "Failed to check password")
		}

		usr.PasswordHash = util.CreateHash(password)

		if err := m.user.Update(projectName, usr); err != nil {
			return errors.Append(err, "Failed to update user password")
		}

		return nil
	})
}

// UserLogout ...
func (m *Manager) UserLogout(projectName string, userID string) *errors.Error {
	if !model.ValidateUserID(userID) {
		return errors.Append(model.ErrUserValidateFailed, "invalid user id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.session.DeleteAll(projectName, userID); err != nil {
			return errors.Append(err, "Delete user session failed")
		}

		return nil
	})
}

// LoginSessionAdd ...
func (m *Manager) LoginSessionAdd(projectName string, ent *model.LoginSession) *errors.Error {
	// create session is in internal only, so validation is not required
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add login session")
		}
		return nil
	})
}

// LoginSessionUpdate ...
func (m *Manager) LoginSessionUpdate(projectName string, ent *model.LoginSession) *errors.Error {
	// update session is in internal only, so validation is not required
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Update(projectName, ent); err != nil {
			return errors.Append(err, "Failed to update login session")
		}
		return nil
	})
}

// LoginSessionDelete ...
func (m *Manager) LoginSessionDelete(projectName string, sessionID string) *errors.Error {
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Delete(projectName, sessionID); err != nil {
			return errors.Append(err, "Failed to delete login session")
		}

		return nil
	})
}

// LoginSessionGet ...
func (m *Manager) LoginSessionGet(projectName string, sessionID string) (*model.LoginSession, *errors.Error) {
	return m.loginSession.Get(projectName, sessionID)
}

// LoginSessionGetByCode ...
func (m *Manager) LoginSessionGetByCode(projectName string, code string) (*model.LoginSession, *errors.Error) {
	return m.loginSession.GetByCode(projectName, code)
}

// SessionAdd ...
func (m *Manager) SessionAdd(projectName string, ent *model.Session) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if _, err := m.session.Get(projectName, ent.SessionID); err != model.ErrNoSuchSession {
			return model.ErrSessionAlreadyExists
		}

		if err := m.session.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add session")
		}
		return nil
	})
}

// SessionGet ..
func (m *Manager) SessionGet(projectName string, sessionID string) (*model.Session, *errors.Error) {
	// TODO(add validation)
	return m.session.Get(projectName, sessionID)
}

// SessionDelete ...
func (m *Manager) SessionDelete(projectName string, sessionID string) *errors.Error {
	if !model.ValidateSessionID(sessionID) {
		return errors.Append(model.ErrSessionValidateFailed, "invalid session id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.session.Delete(projectName, sessionID); err != nil {
			return errors.Append(err, "Failed to revoke session")
		}
		return nil
	})
}

// SessionGetList ...
func (m *Manager) SessionGetList(projectName string, userID string) ([]*model.Session, *errors.Error) {
	if !model.ValidateUserID(userID) {
		return nil, errors.Append(model.ErrSessionValidateFailed, "invalid user id format")
	}

	return m.session.GetList(projectName, userID)
}

// ClientAdd ...
func (m *Manager) ClientAdd(projectName string, ent *model.ClientInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if _, err := m.project.Get(projectName); err != nil {
			return errors.Append(err, "Failed to get project info")
		}

		clis, err := m.client.GetList(ent.ProjectName, &model.ClientFilter{ID: ent.ID})
		if err != nil {
			return errors.Append(err, "Failed to get current client list")
		}

		if len(clis) > 0 {
			return model.ErrClientAlreadyExists
		}

		if err := m.client.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add client")
		}
		return nil
	})
}

// ClientDelete ...
func (m *Manager) ClientDelete(projectName, clientID string) *errors.Error {
	if !model.ValidateClientID(clientID) {
		return errors.Append(model.ErrClientValidateFailed, "invalid client id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		clis, err := m.client.GetList(projectName, &model.ClientFilter{ID: clientID})
		if err != nil {
			return errors.Append(err, "Failed to get current client list")
		}
		if len(clis) == 0 {
			return model.ErrNoSuchClient
		}

		if err := m.loginSession.DeleteAllInClient(projectName, clientID); err != nil {
			return errors.Append(err, "Failed to delete login session of the client")
		}

		if err := m.client.Delete(projectName, clientID); err != nil {
			return errors.Append(err, "Failed to delete client")
		}
		return nil
	})
}

// ClientGetList ...
func (m *Manager) ClientGetList(projectName string, filter *model.ClientFilter) ([]*model.ClientInfo, *errors.Error) {
	if filter != nil && !model.ValidateClientID(filter.ID) {
		return nil, errors.Append(model.ErrClientValidateFailed, "invalid client id format in filter")
	}

	return m.client.GetList(projectName, filter)
}

// ClientGet ...
func (m *Manager) ClientGet(projectName string, clientID string) (*model.ClientInfo, *errors.Error) {
	clients, err := m.ClientGetList(projectName, &model.ClientFilter{ID: clientID})
	if err != nil {
		return nil, err
	}
	if len(clients) == 0 {
		return nil, errors.Append(model.ErrNoSuchClient, "Failed to get client")
	}

	return clients[0], nil
}

// ClientUpdate ...
func (m *Manager) ClientUpdate(projectName string, ent *model.ClientInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		clis, err := m.client.GetList(projectName, &model.ClientFilter{ID: ent.ID})
		if err != nil {
			return errors.Append(err, "Failed to get current client list")
		}
		if len(clis) == 0 {
			return model.ErrNoSuchClient
		}

		if err := m.client.Update(projectName, ent); err != nil {
			return errors.Append(err, "Failed to update client")
		}
		return nil
	})
}

// CustomRoleAdd ...
func (m *Manager) CustomRoleAdd(projectName string, ent *model.CustomRole) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	// TODO(validate name uniquness in project)

	return m.transaction.Transaction(func() *errors.Error {
		_, err := m.customRole.Get(projectName, ent.ID)
		if err != model.ErrNoSuchCustomRole {
			if err == nil {
				return model.ErrCustomRoleAlreadyExists
			}
			return errors.Append(err, "Failed to get customRole info")
		}

		if err := m.customRole.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add customRole")
		}
		return nil
	})
}

// CustomRoleDelete ...
func (m *Manager) CustomRoleDelete(projectName string, customRoleID string) *errors.Error {
	// TODO(validate customRoleID)

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.user.DeleteAllCustomRole(projectName, customRoleID); err != nil {
			return errors.Append(err, "Failed to delete custom role from user")
		}

		if err := m.customRole.Delete(projectName, customRoleID); err != nil {
			return errors.Append(err, "Failed to delete customRole")
		}
		return nil
	})
}

// CustomRoleGetList ...
func (m *Manager) CustomRoleGetList(projectName string, filter *model.CustomRoleFilter) ([]*model.CustomRole, *errors.Error) {
	return m.customRole.GetList(projectName, filter)
}

// CustomRoleGet ...
func (m *Manager) CustomRoleGet(projectName string, customRoleID string) (*model.CustomRole, *errors.Error) {
	// TODO(validate customRoleID)
	return m.customRole.Get(projectName, customRoleID)
}

// CustomRoleUpdate ...
func (m *Manager) CustomRoleUpdate(projectName string, ent *model.CustomRole) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		r, err := m.customRole.GetList(ent.ProjectName, &model.CustomRoleFilter{Name: ent.Name})
		if err != nil {
			return errors.Append(err, "Failed to get role list")
		}

		// check name uniquness in project
		if len(r) > 0 && r[0].ID != ent.ID {
			return model.ErrCustomRoleAlreadyExists
		}

		if err := m.customRole.Update(projectName, ent); err != nil {
			return errors.Append(err, "Failed to update client")
		}
		return nil
	})
}
