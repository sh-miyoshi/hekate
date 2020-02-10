package memory

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"sync"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	// userList[userID] = UserInfo
	userList       map[string]*model.UserInfo
	projectHandler *ProjectInfoHandler
	mu             sync.Mutex
}

// NewUserHandler ...
func NewUserHandler(projectHandler *ProjectInfoHandler) (*UserInfoHandler, error) {
	res := &UserInfoHandler{
		userList:       make(map[string]*model.UserInfo),
		projectHandler: projectHandler,
	}
	return res, nil
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	h.userList[ent.ID] = ent
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(userID string) error {
	if _, exists := h.userList[userID]; exists {
		delete(h.userList, userID)
		return nil
	}
	return model.ErrNoSuchUser
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	res := []string{}

	for _, user := range h.userList {
		if user.ProjectName == projectName {
			res = append(res, user.ID)
		}
	}

	return res, nil
}

// Get ...
func (h *UserInfoHandler) Get(userID string) (*model.UserInfo, error) {
	res, exists := h.userList[userID]
	if !exists {
		return nil, model.ErrNoSuchUser
	}

	return res, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	if _, exists := h.userList[ent.ID]; !exists {
		return model.ErrNoSuchUser
	}

	h.userList[ent.ID] = ent

	return nil
}

// GetByName ...
func (h *UserInfoHandler) GetByName(projectName string, userName string) (*model.UserInfo, error) {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return nil, model.ErrNoSuchProject
	}

	for _, user := range h.userList {
		if user.ProjectName == projectName && user.Name == userName {
			return user, nil
		}
	}
	return nil, model.ErrNoSuchUser
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) error {
	for _, user := range h.userList {
		if user.ProjectName == projectName {
			delete(h.userList, user.ID)
		}
	}
	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(userID string, roleType model.RoleType, roleID string) error {
	if _, exists := h.userList[userID]; !exists {
		return model.ErrNoSuchUser
	}

	roles := h.userList[userID].SystemRoles
	if roleType == model.RoleCustom {
		roles = h.userList[userID].CustomRoles
	}

	for _, r := range roles {
		if r == roleID {
			return model.ErrRoleAlreadyAppended
		}
	}

	roles = append(roles, roleID)
	if roleType == model.RoleCustom {
		h.userList[userID].CustomRoles = roles
	} else if roleType == model.RoleSystem {
		h.userList[userID].SystemRoles = roles
	}

	return nil
}

// DeleteRole ....
func (h *UserInfoHandler) DeleteRole(userID string, roleID string) error {
	if _, exists := h.userList[userID]; !exists {
		return model.ErrNoSuchUser
	}

	deleted := false
	roles := []string{}
	for _, r := range h.userList[userID].SystemRoles {
		if r == roleID {
			deleted = true
		} else {
			roles = append(roles, r)
		}
	}

	if deleted {
		h.userList[userID].SystemRoles = roles
		return nil
	}

	deleted = false
	roles = []string{}
	for _, r := range h.userList[userID].CustomRoles {
		if r == roleID {
			deleted = true
		} else {
			roles = append(roles, r)
		}
	}

	if deleted {
		h.userList[userID].CustomRoles = roles
		return nil
	}

	return model.ErrNoSuchRoleInUser
}

// AddLoginSession ...
func (h *UserInfoHandler) AddLoginSession(userID string, info *model.LoginSessionInfo) error {
	if _, exists := h.userList[userID]; !exists {
		return model.ErrNoSuchUser
	}

	for _, s := range h.userList[userID].LoginSessions {
		if s.VerifyCode == info.VerifyCode {
			return model.ErrLoginSessionAlreadyExists
		}
	}

	h.userList[userID].LoginSessions = append(h.userList[userID].LoginSessions, info)

	return nil
}

// DeleteLoginSession ...
func (h *UserInfoHandler) DeleteLoginSession(userID string, code string) error {
	if _, exists := h.userList[userID]; !exists {
		return model.ErrNoSuchUser
	}

	found := false
	sessions := []*model.LoginSessionInfo{}
	for _, s := range h.userList[userID].LoginSessions {
		if s.VerifyCode == code {
			found = true
		} else {
			sessions = append(sessions, s)
		}
	}

	if !found {
		return model.ErrNoSuchLoginSession
	}

	h.userList[userID].LoginSessions = sessions
	return nil
}

// BeginTx ...
func (h *UserInfoHandler) BeginTx() error {
	h.mu.Lock()
	return nil
}

// CommitTx ...
func (h *UserInfoHandler) CommitTx() error {
	h.mu.Unlock()
	return nil
}

// AbortTx ...
func (h *UserInfoHandler) AbortTx() error {
	h.mu.Unlock()
	return nil
}
