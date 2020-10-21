package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	// userList[userID] = UserInfo
	userList map[string]*model.UserInfo
}

// NewUserHandler ...
func NewUserHandler() *UserInfoHandler {
	res := &UserInfoHandler{
		userList: make(map[string]*model.UserInfo),
	}
	return res
}

// Add ...
func (h *UserInfoHandler) Add(projectName string, ent *model.UserInfo) *errors.Error {
	h.userList[ent.ID] = ent
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) *errors.Error {
	if res, exists := h.userList[userID]; exists {
		if res.ProjectName == projectName {
			delete(h.userList, userID)
			return nil
		}
	}
	return model.ErrNoSuchUser
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, *errors.Error) {
	res := []*model.UserInfo{}

	for _, user := range h.userList {
		if user.ProjectName == projectName {
			res = append(res, user)
		}
	}

	if filter != nil {
		res = matchFilterUserList(res, projectName, filter)
	}

	return res, nil
}

// Update ...
func (h *UserInfoHandler) Update(projectName string, ent *model.UserInfo) *errors.Error {
	if res, exists := h.userList[ent.ID]; !exists || res.ProjectName != projectName {
		return model.ErrNoSuchUser
	}

	h.userList[ent.ID] = ent

	return nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) *errors.Error {
	for _, user := range h.userList {
		if user.ProjectName == projectName {
			delete(h.userList, user.ID)
		}
	}
	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(projectName string, userID string, roleType model.RoleType, roleID string) *errors.Error {
	if res, exists := h.userList[userID]; !exists || res.ProjectName != projectName {
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
func (h *UserInfoHandler) DeleteRole(projectName string, userID string, roleID string) *errors.Error {
	if res, exists := h.userList[userID]; !exists || res.ProjectName != projectName {
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

// DeleteAllCustomRole ...
func (h *UserInfoHandler) DeleteAllCustomRole(projectName string, roleID string) *errors.Error {
	for id, user := range h.userList {
		if user.ProjectName != projectName {
			continue
		}

		deleted := false
		roles := []string{}

		for _, r := range user.CustomRoles {
			if roleID == r {
				deleted = true
			} else {
				roles = append(roles, r)
			}
		}

		if deleted {
			h.userList[id].CustomRoles = roles
		}
	}
	return nil
}

func matchFilterUserList(data []*model.UserInfo, projectName string, filter *model.UserFilter) []*model.UserInfo {
	if filter == nil {
		return data
	}
	res := []*model.UserInfo{}

	for _, user := range data {
		if projectName == user.ProjectName {
			if filter.Name != "" && user.Name != filter.Name {
				// missmatch name
				continue
			}
			if filter.ID != "" && user.ID != filter.ID {
				// missmatch id
				continue
			}
		}
		res = append(res, user)
	}

	return res
}
