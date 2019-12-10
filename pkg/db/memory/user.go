package memory

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	// userList[projectName][userID] = UserInfo
	userList       map[string](map[string]model.UserInfo)
	projectHandler *ProjectInfoHandler
}

// NewUserHandler ...
func NewUserHandler(projectHandler *ProjectInfoHandler) (*UserInfoHandler, error) {
	res := &UserInfoHandler{
		userList:       make(map[string](map[string]model.UserInfo)),
		projectHandler: projectHandler,
	}
	return res, nil
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	if err := ent.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate entry")
	}

	if _, err := h.projectHandler.Get(ent.ProjectName); err != nil {
		return errors.Wrap(err, "Failed to get project")
	}

	// If userList do not contains project info, create project info
	if _, exists := h.userList[ent.ProjectName]; !exists {
		h.userList[ent.ProjectName] = make(map[string]model.UserInfo)
	}

	if _, exists := h.userList[ent.ProjectName][ent.ID]; exists {
		return errors.Cause(model.ErrUserAlreadyExists)
	}

	for _, user := range h.userList[ent.ProjectName] {
		if user.Name == ent.Name {
			return errors.Cause(model.ErrUserAlreadyExists)
		}
	}

	h.userList[ent.ProjectName][ent.ID] = *ent
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) error {
	if _, exists := h.userList[projectName]; !exists {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[projectName][userID]; exists {
		delete(h.userList[projectName], userID)
		return nil
	}
	return errors.Cause(model.ErrNoSuchUser)
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	res := []string{}

	if _, exists := h.userList[projectName]; !exists {
		// project is created in Add method, so maybe empty project
		return res, nil
	}

	for _, user := range h.userList[projectName] {
		res = append(res, user.ID)
	}

	return res, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectName string, userID string) (*model.UserInfo, error) {
	if _, exists := h.userList[projectName]; !exists {
		return nil, errors.Cause(model.ErrNoSuchProject)
	}

	res, exists := h.userList[projectName][userID]
	if !exists {
		return nil, errors.Cause(model.ErrNoSuchUser)
	}

	return &res, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	if _, exists := h.userList[ent.ProjectName]; !exists {
		return errors.Cause(model.ErrNoSuchProject)
	}

	if _, exists := h.userList[ent.ProjectName][ent.ID]; !exists {
		return errors.Cause(model.ErrNoSuchUser)
	}

	h.userList[ent.ProjectName][ent.ID] = *ent

	return nil
}

// GetIDByName ...
func (h *UserInfoHandler) GetIDByName(projectName string, userName string) (string, error) {
	if _, exists := h.userList[projectName]; !exists {
		return "", errors.Cause(model.ErrNoSuchProject)
	}

	for _, user := range h.userList[projectName] {
		if user.Name == userName {
			return user.ID, nil
		}
	}
	return "", errors.Cause(model.ErrNoSuchUser)
}

// DeleteProjectDefine ...
func (h *UserInfoHandler) DeleteProjectDefine(projectName string) error {
	if _, exists := h.userList[projectName]; exists {
		delete(h.userList, projectName)
	}
	return nil
}
