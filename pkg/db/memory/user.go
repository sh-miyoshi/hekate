package memory

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	// userList[projectID][userID] = UserInfo
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

	if _, err := h.projectHandler.Get(ent.ProjectID); err != nil {
		return errors.Wrap(err, "Failed to get project")
	}

	// If userList do not contains project info, create project info
	if _, exists := h.userList[ent.ProjectID]; !exists {
		h.userList[ent.ProjectID] = make(map[string]model.UserInfo)
	}

	if _, exists := h.userList[ent.ProjectID][ent.ID]; exists {
		return errors.Cause(model.ErrUserAlreadyExists)
	}

	for _, user := range h.userList[ent.ProjectID] {
		if user.Name == ent.Name {
			return errors.Cause(model.ErrUserAlreadyExists)
		}
	}

	h.userList[ent.ProjectID][ent.ID] = *ent
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectID string, userID string) error {
	// TODO(not implemented yet)
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectID string) ([]string, error) {
	// TODO(not implemented yet)
	return []string{}, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectID string, userID string) (*model.UserInfo, error) {
	if _, exists := h.userList[projectID]; !exists {
		return nil, errors.Cause(model.ErrNoSuchProject)
	}

	res, exists := h.userList[projectID][userID]
	if !exists {
		return nil, errors.Cause(model.ErrNoSuchUser)
	}

	return &res, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	// TODO(not implemented yet)
	return nil
}

// GetIDByName ...
func (h *UserInfoHandler) GetIDByName(projectID string, userName string) (string, error) {
	if _, exists := h.userList[projectID]; !exists {
		return "", errors.Cause(model.ErrNoSuchProject)
	}

	for _, user := range h.userList[projectID] {
		if user.Name == userName {
			return user.ID, nil
		}
	}
	return "", errors.Cause(model.ErrNoSuchUser)
}

// DeleteProjectDefine ...
func (h *UserInfoHandler) DeleteProjectDefine(projectID string) error {
	// TODO(not implemented yet)
	return nil
}
