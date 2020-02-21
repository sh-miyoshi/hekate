package memory

import (
	"sync"

	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
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
func (h *UserInfoHandler) GetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, error) {
	res := []*model.UserInfo{}

	for _, user := range h.userList {
		if user.ProjectName == projectName {
			res = append(res, user)
		}
	}

	// TODO
	if filter != nil {
		res = filterUserList(res, filter)
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

func filterUserList(data []*model.UserInfo, filter *model.UserFilter) []*model.UserInfo {
	if filter == nil {
		return data
	}
	res := []*model.UserInfo{}

	for _, user := range data {
		if filter.Name != "" && user.Name != filter.Name {
			// missmatch name
			continue
		}
		// TODO(add other filter)
		res = append(res, user)
	}

	return res
}
