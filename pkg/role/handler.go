package role

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

// Handler ...
type Handler struct {
	roleList []Info
}

var inst *Handler

// InitHandler ...
func InitHandler() error {
	if inst != nil {
		return errors.New("Default Role Handler is already initialized")
	}

	inst = &Handler{}

	// Create default role
	inst.createRole(ResCluster, TypeRead)
	inst.createRole(ResCluster, TypeWrite)
	inst.createRole(ResProject, TypeRead)
	inst.createRole(ResProject, TypeWrite)
	inst.createRole(ResRole, TypeRead)
	inst.createRole(ResRole, TypeWrite)
	inst.createRole(ResUser, TypeRead)
	inst.createRole(ResUser, TypeWrite)
	inst.createRole(ResClient, TypeRead)
	inst.createRole(ResClient, TypeWrite)
	inst.createRole(ResCustomRole, TypeRead)
	inst.createRole(ResCustomRole, TypeWrite)

	roles := []string{}
	for _, role := range inst.roleList {
		roles = append(roles, role.Name)
	}
	logger.Debug("All Default Role List: %v", roles)

	return nil
}

// Authorize ...
func Authorize(roles []string, targetResource Resource, roleType Type) bool {
	name := fmt.Sprintf("%s-%s", roleType.String(), targetResource.String())
	logger.Debug("Auth want: %s, have: %v", name, roles)

	for _, role := range roles {
		if role == name {
			return true
		}
	}
	return false
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

// Parse ...
func (h *Handler) Parse(role string) (*Resource, *Type, bool) {
	for _, r := range h.roleList {
		if r.ID == role {
			return &r.TargetResource, &r.RoleType, true
		}
	}
	return nil, nil, false
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
