package db

import (
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/sh-miyoshi/hekate/pkg/db/memory"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/secret"
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
	device       model.DeviceHandler

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
			device:       memory.NewDeviceHandler(),
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
		deviceHandler, err := mongo.NewDeviceHandler(dbClient)
		if err != nil {
			return errors.Append(err, "Failed to create device handler")
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
			device:       deviceHandler,
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

	keys, err := secret.GetSignKey(ent.TokenConfig.SigningAlgorithm)
	if err != nil {
		return err
	}
	ent.TokenConfig.SignSecretKey = keys.Private
	ent.TokenConfig.SignPublicKey = keys.Public

	return m.transaction.Transaction(func() *errors.Error {
		prjs, err := m.project.GetList(&model.ProjectFilter{Name: ent.Name})
		if err != nil {
			return errors.Append(err, "Failed to get current project list")
		}

		if len(prjs) != 0 {
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
		prjs, err := m.project.GetList(&model.ProjectFilter{Name: name})
		if err != nil {
			return errors.Append(err, "Failed to get delete project info")
		}
		if len(prjs) == 0 {
			return model.ErrNoSuchProject
		}

		prj := prjs[0]

		if !prj.PermitDelete {
			return errors.Append(model.ErrDeleteBlockedProject, "the project can not delete")
		}

		if err := m.loginSession.DeleteAll(name); err != nil {
			return errors.Append(err, "Failed to delete login session data")
		}

		if err := m.session.DeleteAll(name); err != nil {
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

		if err := m.device.DeleteAll(name); err != nil {
			return errors.Append(err, "Failed to delete device data")
		}

		if err := m.project.Delete(name); err != nil {
			return errors.Append(err, "Failed to delete project")
		}

		return nil
	})
}

// ProjectGetList ...
func (m *Manager) ProjectGetList(filter *model.ProjectFilter) ([]*model.ProjectInfo, *errors.Error) {
	if filter != nil {
		if filter.Name != "" && !model.ValidateProjectName(filter.Name) {
			return nil, errors.Append(model.ErrProjectValidateFailed, "Invalid project name format")
		}
	}
	return m.project.GetList(filter)
}

// ProjectGet ...
func (m *Manager) ProjectGet(name string) (*model.ProjectInfo, *errors.Error) {
	prjs, err := m.ProjectGetList(&model.ProjectFilter{Name: name})
	if err != nil {
		return nil, err
	}

	if len(prjs) == 0 {
		return nil, model.ErrNoSuchProject
	}

	return prjs[0], nil
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

// ProjectSecretReset ...
func (m *Manager) ProjectSecretReset(name string) *errors.Error {
	prj, err := m.ProjectGet(name)
	if err != nil {
		return err
	}

	keys, err := secret.GetSignKey(prj.TokenConfig.SigningAlgorithm)
	if err != nil {
		return err
	}
	prj.TokenConfig.SignSecretKey = keys.Private
	prj.TokenConfig.SignPublicKey = keys.Public

	return m.ProjectUpdate(prj)
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
			roles, err := m.customRole.GetList(projectName, &model.CustomRoleFilter{ID: r})
			if err != nil {
				return errors.Append(err, "Custom role get error")
			}
			if len(roles) == 0 {
				return errors.Append(model.ErrUserValidateFailed, "No such custom role")
			}
		}

		users, err := m.user.GetList(projectName, &model.UserFilter{ID: ent.ID})
		if err != nil {
			return errors.Append(err, "Failed to get user info")
		}
		if len(users) != 0 {
			return model.ErrUserAlreadyExists
		}

		// Check duplicate user by name
		users, err = m.user.GetList(ent.ProjectName, &model.UserFilter{Name: ent.Name})
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
		if _, err := m.UserGet(projectName, userID); err != nil {
			return errors.Append(err, "Failed to get user info")
		}

		if err := m.loginSession.Delete(projectName, &model.LoginSessionFilter{UserID: userID}); err != nil {
			return errors.Append(err, "Delete authoriation code failed")
		}

		if err := m.session.Delete(projectName, &model.SessionFilter{UserID: userID}); err != nil {
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
	if filter != nil {
		if filter.ID != "" && !model.ValidateUserID(filter.ID) {
			return nil, errors.Append(model.ErrUserValidateFailed, "Invalid user id format")
		}
		if filter.Name != "" && !model.ValidateUserName(filter.Name) {
			return nil, errors.Append(model.ErrUserValidateFailed, "Invalid user name format")
		}
	}

	return m.user.GetList(projectName, filter)
}

// UserGet ...
func (m *Manager) UserGet(projectName string, userID string) (*model.UserInfo, *errors.Error) {
	users, err := m.UserGetList(projectName, &model.UserFilter{ID: userID})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, model.ErrNoSuchUser
	}

	return users[0], nil
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
			roles, err := m.customRole.GetList(projectName, &model.CustomRoleFilter{ID: r})
			if err != nil {
				return errors.Append(err, "Custom role get error")
			}
			if len(roles) == 0 {
				return errors.Append(model.ErrUserValidateFailed, "No such custom role")
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
		users, err := m.user.GetList(projectName, &model.UserFilter{ID: userID})
		if err != nil {
			return errors.Append(err, "Failed to get current user list")
		}
		if len(users) == 0 {
			return model.ErrNoSuchUser
		}

		// Validate RoleID
		if roleType == model.RoleSystem {
			res, typ, ok := role.GetInst().Parse(roleID)
			if !ok {
				return errors.Append(model.ErrUserValidateFailed, "Invalid system role")
			}

			usr := users[0]
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
			roles, err := m.customRole.GetList(projectName, &model.CustomRoleFilter{ID: roleID})
			if err != nil {
				return errors.Append(err, "Custom role get error")
			}
			if len(roles) == 0 {
				return errors.Append(model.ErrUserValidateFailed, "No such custom role")
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
		users, err := m.user.GetList(projectName, &model.UserFilter{ID: userID})
		if err != nil {
			return errors.Append(err, "Failed to get current user list")
		}
		if len(users) == 0 {
			return model.ErrNoSuchUser
		}
		usr := users[0]

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
		prjs, err := m.project.GetList(&model.ProjectFilter{Name: projectName})
		if err != nil {
			return errors.Append(err, "Failed to get project associated with the user")
		}
		if len(prjs) == 0 {
			return model.ErrNoSuchUser
		}
		prj := prjs[0]

		users, err := m.user.GetList(projectName, &model.UserFilter{ID: userID})
		if err != nil {
			return errors.Append(err, "Failed to get user of change password")
		}
		if len(users) == 0 {
			return model.ErrNoSuchUser
		}
		usr := users[0]

		if err := secret.CheckPassword(usr.Name, password, prj.PasswordPolicy); err != nil {
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
		if err := m.session.Delete(projectName, &model.SessionFilter{UserID: userID}); err != nil {
			return errors.Append(err, "Delete user session failed")
		}

		return nil
	})
}

// LoginSessionAdd ...
func (m *Manager) LoginSessionAdd(projectName string, ent *model.LoginSession) *errors.Error {
	// login session is in internal only, so validation is not required
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add login session")
		}
		return nil
	})
}

// LoginSessionUpdate ...
func (m *Manager) LoginSessionUpdate(projectName string, ent *model.LoginSession) *errors.Error {
	// login session is in internal only, so validation is not required
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Update(projectName, ent); err != nil {
			return errors.Append(err, "Failed to update login session")
		}
		return nil
	})
}

// LoginSessionDelete ...
func (m *Manager) LoginSessionDelete(projectName string, sessionID string) *errors.Error {
	// login session is in internal only, so validation is not required
	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Delete(projectName, &model.LoginSessionFilter{SessionID: sessionID}); err != nil {
			return errors.Append(err, "Failed to delete login session")
		}

		return nil
	})
}

// LoginSessionGet ...
func (m *Manager) LoginSessionGet(projectName string, sessionID string) (*model.LoginSession, *errors.Error) {
	// login session is in internal only, so validation is not required
	return m.loginSession.Get(projectName, sessionID)
}

// LoginSessionGetByCode ...
func (m *Manager) LoginSessionGetByCode(projectName string, code string) (*model.LoginSession, *errors.Error) {
	// login session is in internal only, so validation is not required
	return m.loginSession.GetByCode(projectName, code)
}

// SessionAdd ...
func (m *Manager) SessionAdd(projectName string, ent *model.Session) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		sessions, err := m.session.GetList(projectName, &model.SessionFilter{SessionID: ent.SessionID})
		if err != nil {
			return errors.Append(err, "Failed to get current session list")
		}
		if len(sessions) != 0 {
			return model.ErrSessionAlreadyExists
		}

		if err := m.session.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add session")
		}
		return nil
	})
}

// SessionGetList ...
func (m *Manager) SessionGetList(projectName string, filter *model.SessionFilter) ([]*model.Session, *errors.Error) {
	if filter != nil {
		if filter.SessionID != "" && !model.ValidateSessionID(filter.SessionID) {
			return nil, model.ErrSessionValidateFailed
		}
		if filter.UserID != "" && !model.ValidateUserID(filter.UserID) {
			return nil, model.ErrSessionValidateFailed
		}
	}

	return m.session.GetList(projectName, filter)
}

// SessionGet ..
func (m *Manager) SessionGet(projectName string, sessionID string) (*model.Session, *errors.Error) {
	sessions, err := m.SessionGetList(projectName, &model.SessionFilter{SessionID: sessionID})
	if err != nil {
		return nil, errors.Append(err, "Failed to get current session list")
	}
	if len(sessions) == 0 {
		return nil, model.ErrNoSuchSession
	}
	return sessions[0], nil
}

// SessionDelete ...
func (m *Manager) SessionDelete(projectName string, sessionID string) *errors.Error {
	if !model.ValidateSessionID(sessionID) {
		return errors.Append(model.ErrSessionValidateFailed, "invalid session id format")
	}

	return m.transaction.Transaction(func() *errors.Error {
		sessions, err := m.SessionGetList(projectName, &model.SessionFilter{SessionID: sessionID})
		if err != nil {
			return errors.Append(err, "Failed to get current session list")
		}

		if len(sessions) == 0 {
			return model.ErrNoSuchSession
		}

		if err := m.session.Delete(projectName, &model.SessionFilter{SessionID: sessionID}); err != nil {
			return errors.Append(err, "Failed to revoke session")
		}
		return nil
	})
}

// ClientAdd ...
func (m *Manager) ClientAdd(projectName string, ent *model.ClientInfo) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		prjs, err := m.project.GetList(&model.ProjectFilter{Name: projectName})
		if err != nil {
			return errors.Append(err, "Failed to get current project")
		}
		if len(prjs) == 0 {
			return model.ErrNoSuchProject
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

		if err := m.loginSession.Delete(projectName, &model.LoginSessionFilter{ClientID: clientID}); err != nil {
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

	return m.transaction.Transaction(func() *errors.Error {
		prjs, err := m.project.GetList(&model.ProjectFilter{Name: projectName})
		if err != nil {
			return errors.Append(err, "Failed to get current project")
		}
		if len(prjs) == 0 {
			return model.ErrNoSuchProject
		}

		roles, err := m.customRole.GetList(projectName, &model.CustomRoleFilter{
			ID:   ent.ID,
			Name: ent.Name,
		})
		if err != nil {
			return errors.Append(err, "Failed to get current custom role list")
		}
		if len(roles) != 0 {
			return model.ErrCustomRoleAlreadyExists
		}

		if err := m.customRole.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add customRole")
		}
		return nil
	})
}

// CustomRoleDelete ...
func (m *Manager) CustomRoleDelete(projectName string, customRoleID string) *errors.Error {
	if !model.ValidateCustomRoleID(customRoleID) {
		return model.ErrCustomRoleValidateFailed
	}

	return m.transaction.Transaction(func() *errors.Error {
		roles, err := m.customRole.GetList(projectName, &model.CustomRoleFilter{ID: customRoleID})
		if err != nil {
			return errors.Append(err, "Failed to get current custom role list")
		}
		if len(roles) == 0 {
			return model.ErrNoSuchCustomRole
		}

		// Delete custom role from all user
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
	if filter != nil {
		if filter.ID != "" && !model.ValidateCustomRoleID(filter.ID) {
			return nil, errors.Append(model.ErrCustomRoleValidateFailed, "Invalid role id format")
		}
		if filter.Name != "" && !model.ValidateCustomRoleName(filter.Name) {
			return nil, errors.Append(model.ErrCustomRoleValidateFailed, "Invalid role name format")
		}
	}
	return m.customRole.GetList(projectName, filter)
}

// CustomRoleGet ...
func (m *Manager) CustomRoleGet(projectName string, customRoleID string) (*model.CustomRole, *errors.Error) {
	roles, err := m.CustomRoleGetList(projectName, &model.CustomRoleFilter{ID: customRoleID})
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.Append(model.ErrNoSuchCustomRole, "Failed to get role")
	}

	return roles[0], nil
}

// CustomRoleUpdate ...
func (m *Manager) CustomRoleUpdate(projectName string, ent *model.CustomRole) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		r, err := m.customRole.GetList(ent.ProjectName, &model.CustomRoleFilter{ID: ent.ID})
		if err != nil {
			return errors.Append(err, "Failed to get current role list by ID")
		}

		if len(r) == 0 {
			return model.ErrNoSuchCustomRole
		}

		r, err = m.customRole.GetList(ent.ProjectName, &model.CustomRoleFilter{Name: ent.Name})
		if err != nil {
			return errors.Append(err, "Failed to get current role list by Name")
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

// DeviceAdd ...
func (m *Manager) DeviceAdd(projectName string, ent *model.Device) *errors.Error {
	if err := ent.Validate(); err != nil {
		return errors.Append(err, "Failed to validate entry")
	}

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.device.Add(projectName, ent); err != nil {
			return errors.Append(err, "Failed to add device")
		}
		return nil
	})
}

// DeviceDelete ...
func (m *Manager) DeviceDelete(projectName string, deviceCode string) *errors.Error {
	// TODO validate deviceCode

	return m.transaction.Transaction(func() *errors.Error {
		devices, err := m.device.GetList(projectName, &model.DeviceFilter{DeviceCode: deviceCode})
		if err != nil {
			return errors.Append(err, "Failed to get current device list")
		}
		if len(devices) == 0 {
			return model.ErrNoSuchDevice
		}

		if err := m.loginSession.Delete(projectName, &model.LoginSessionFilter{SessionID: devices[0].LoginSessionID}); err != nil {
			return errors.Append(err, "Failed to delete login session")
		}
		if err := m.device.Delete(projectName, deviceCode); err != nil {
			return errors.Append(err, "Failed to delete device")
		}

		return nil
	})
}

// DeviceGetList ...
func (m *Manager) DeviceGetList(projectName string, filter *model.DeviceFilter) ([]*model.Device, *errors.Error) {
	// TODO validate filter
	return m.device.GetList(projectName, filter)
}

// DeleteExpiredSessions ...
func (m *Manager) DeleteExpiredSessions() *errors.Error {
	now := time.Now()

	return m.transaction.Transaction(func() *errors.Error {
		if err := m.loginSession.Cleanup(now); err != nil {
			return errors.Append(err, "Failed to cleanup login sessions")
		}

		if err := m.session.Cleanup(now); err != nil {
			return errors.Append(err, "Failed to cleanup sessions")
		}

		if err := m.device.Cleanup(now); err != nil {
			return errors.Append(err, "Failed to cleanup devices")
		}

		return nil
	})
}
