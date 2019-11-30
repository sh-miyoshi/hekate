package defaultrole

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// Handler ...
type Handler struct {
	roleList []RoleInfo
}

var inst *Handler

// InitHandler ...
func InitHandler() error {
	if inst != nil {
		return errors.Cause(fmt.Errorf("Default Role Handler is already initialized"))
	}

	resources := []string{
		"project", "role", "user", "cluster",
	}
	types := []string{
		"read", "write", "manage",
	}

	inst = &Handler{}

	// Create default role
	for _, res := range resources {
		for _, typ := range types {
			inst.createRole(res, typ)
		}
	}

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

// GetByID ...
func (h *Handler) GetByID(id string) (*RoleInfo, error) {
	for _, role := range h.roleList {
		if role.ID == id {
			return &role, nil
		}
	}
	return nil, errors.New("No such role")
}

func (h *Handler) createRole(targetResource string, roleType string) {
	val := RoleInfo{
		ID:             uuid.New().String(),
		Name:           fmt.Sprintf("%s-%s", roleType, targetResource),
		TargetResource: targetResource,
		RoleType:       roleType,
	}
	h.roleList = append(h.roleList, val)
}
