package memory

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	// userList[userID] = UserInfo
	userList       map[string]*model.UserInfo
	projectHandler *ProjectInfoHandler
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
func (h *UserInfoHandler) Delete(projectName string, userID string) error {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[userID]; exists {
		delete(h.userList, userID)
		return nil
	}
	return errors.Cause(model.ErrNoSuchUser)
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	res := []string{}

	if _, err := h.projectHandler.Get(projectName); err != nil {
		// project is created in Add method, so maybe empty project
		return res, nil
	}

	for _, user := range h.userList {
		if user.ProjectName == projectName {
			res = append(res, user.ID)
		}
	}

	return res, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectName string, userID string) (*model.UserInfo, error) {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return nil, errors.Cause(model.ErrNoSuchProject)
	}

	res, exists := h.userList[userID]
	if !exists {
		return nil, errors.Cause(model.ErrNoSuchUser)
	}

	return res, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	if _, err := h.projectHandler.Get(ent.ProjectName); err != nil {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[ent.ID]; !exists {
		return errors.Cause(model.ErrNoSuchUser)
	}

	h.userList[ent.ID] = ent

	return nil
}

// GetByName ...
func (h *UserInfoHandler) GetByName(projectName string, userName string) (*model.UserInfo, error) {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return nil, errors.Cause(model.ErrNoSuchProject)
	}

	for _, user := range h.userList {
		if user.ProjectName == projectName && user.Name == userName {
			return user, nil
		}
	}
	return nil, errors.Cause(model.ErrNoSuchUser)
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
func (h *UserInfoHandler) AddRole(projectName string, userID string, roleID string) error {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[userID]; !exists {
		return errors.Cause(model.ErrNoSuchUser)
	}

	roles := h.userList[userID].Roles
	for _, r := range roles {
		if r == roleID {
			return errors.Cause(model.ErrRoleAlreadyAppended)
		}
	}

	roles = append(roles, roleID)
	h.userList[userID].Roles = roles

	return nil
}

// DeleteRole ....
func (h *UserInfoHandler) DeleteRole(projectName string, userID string, roleID string) error {
	if _, err := h.projectHandler.Get(projectName); err != nil {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[userID]; !exists {
		return errors.Cause(model.ErrNoSuchUser)
	}

	deleted := false
	roles := []string{}
	for _, r := range h.userList[userID].Roles {
		if r == roleID {
			deleted = true
		} else {
			roles = append(roles, r)
		}
	}

	h.userList[userID].Roles = roles

	if !deleted {
		return errors.Cause(model.ErrNoSuchRoleInUser)
	}

	return nil
}
