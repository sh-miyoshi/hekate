package output

import (
	"encoding/json"
	"fmt"

	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
)

// CustomRoleFormat ...
type CustomRoleFormat struct {
	role *roleapi.CustomRoleGetResponse
}

// CustomRolesFormat ...
type CustomRolesFormat struct {
	roles []*roleapi.CustomRoleGetResponse
}

// NewCustomRoleFormat ...
func NewCustomRoleFormat(role *roleapi.CustomRoleGetResponse) *CustomRoleFormat {
	return &CustomRoleFormat{
		role: role,
	}
}

// NewRolesInfoFormat ...
func NewRolesInfoFormat(roles []*roleapi.CustomRoleGetResponse) *CustomRolesFormat {
	return &CustomRolesFormat{
		roles: roles,
	}
}

// ToText ...
func (f *CustomRoleFormat) ToText() (string, error) {
	res := fmt.Sprintf("ID:           %s\n", f.role.ID)
	res += fmt.Sprintf("Name:         %s\n", f.role.Name)
	res += fmt.Sprintf("Created Time: %s\n", f.role.CreatedAt)
	return res, nil
}

// ToJSON ...
func (f *CustomRoleFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.role)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToText ...
func (f *CustomRolesFormat) ToText() (string, error) {
	res := ""
	for i, prj := range f.roles {
		format := NewCustomRoleFormat(prj)
		msg, err := format.ToText()
		if err != nil {
			return "", err
		}
		res += msg
		if i < len(f.roles)-1 {
			res += "\n---\n"
		}
	}
	return res, nil
}

// ToJSON ...
func (f *CustomRolesFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.roles)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
