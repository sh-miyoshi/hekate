package role

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Handler ...
type Handler struct {
	roleList []Info
}

var inst *Handler

// InitHandler ...
func InitHandler() error {
	if inst != nil {
		return errors.Cause(fmt.Errorf("Default Role Handler is already initialized"))
	}

	inst = &Handler{}

	// Create default role
	inst.createRole(ResCluster, TypeManage)
	inst.createRole(ResCluster, TypeRead)
	inst.createRole(ResCluster, TypeWrite)
	inst.createRole(ResProject, TypeManage)
	inst.createRole(ResProject, TypeRead)
	inst.createRole(ResProject, TypeWrite)
	inst.createRole(ResRole, TypeManage)
	inst.createRole(ResRole, TypeRead)
	inst.createRole(ResRole, TypeWrite)
	inst.createRole(ResUser, TypeManage)
	inst.createRole(ResUser, TypeRead)
	inst.createRole(ResUser, TypeWrite)

	roles := []string{}
	for _, role := range inst.roleList {
		roles = append(roles, role.Name)
	}
	logger.Debug("All Default Role List: %v", roles)

	return nil
}

// GetInst returns an instance of DB Manager
func GetInst() *Handler {
	return inst
}

// GetList ...
func (h *Handler) GetList() []string {
	res := []string{}
	for _, role := range h.roleList {
		res = append(res, role.ID)
	}
	return res
}

// Authorize ...
func (h *Handler) Authorize(roles []string, targetResource Resource, roleType Type) bool {
	name := fmt.Sprintf("%s-%s", roleType.String(), targetResource.String())
	logger.Debug("Auth want: %s, have: %v", name, roles)

	for _, role := range roles {
		if role == name {
			return true
		}
	}
	return false
}

func (h *Handler) createRole(targetResource Resource, roleType Type) {
	name := fmt.Sprintf("%s-%s", roleType.String(), targetResource.String())
	val := Info{
		ID:             name,
		Name:           name,
		TargetResource: targetResource,
		RoleType:       roleType,
	}
	h.roleList = append(h.roleList, val)
}
