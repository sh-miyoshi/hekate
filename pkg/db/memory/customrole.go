package memory

import (
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

// CustomRoleHandler implement db.CustomRoleHandler
type CustomRoleHandler struct {
	// roleList[roleID] = CustomRole
	roleList map[string]*model.CustomRole
}

// NewCustomRoleHandler ...
func NewCustomRoleHandler() *CustomRoleHandler {
	res := &CustomRoleHandler{
		roleList: make(map[string]*model.CustomRole),
	}
	return res
}

// Add ...
func (h *CustomRoleHandler) Add(ent *model.CustomRole) error {
	h.roleList[ent.ID] = ent
	return nil
}

// Delete ...
func (h *CustomRoleHandler) Delete(roleID string) error {
	if _, exists := h.roleList[roleID]; exists {
		delete(h.roleList, roleID)
		return nil
	}
	return model.ErrNoSuchCustomRole
}

// GetList ...
func (h *CustomRoleHandler) GetList(projectName string, filter *model.CustomRoleFilter) ([]*model.CustomRole, error) {
	res := []*model.CustomRole{}

	for _, role := range h.roleList {
		if role.ProjectName == projectName {
			res = append(res, role)
		}
	}

	if filter != nil {
		res = filterRoleList(res, filter)
	}

	return res, nil
}

// Get ...
func (h *CustomRoleHandler) Get(roleID string) (*model.CustomRole, error) {
	res, exists := h.roleList[roleID]
	if !exists {
		return nil, model.ErrNoSuchCustomRole
	}

	return res, nil
}

// Update ...
func (h *CustomRoleHandler) Update(ent *model.CustomRole) error {
	if _, exists := h.roleList[ent.ID]; !exists {
		return model.ErrNoSuchCustomRole
	}

	h.roleList[ent.ID] = ent

	return nil
}

// DeleteAll ...
func (h *CustomRoleHandler) DeleteAll(projectName string) error {
	for _, role := range h.roleList {
		if role.ProjectName == projectName {
			delete(h.roleList, role.ID)
		}
	}
	return nil
}

func filterRoleList(data []*model.CustomRole, filter *model.CustomRoleFilter) []*model.CustomRole {
	if filter == nil {
		return data
	}
	res := []*model.CustomRole{}

	for _, role := range data {
		if filter.Name != "" && role.Name != filter.Name {
			// missmatch name
			continue
		}
		// TODO(add other filter)
		res = append(res, role)
	}

	return res
}
