package output

import (
	"encoding/json"
	"fmt"

	userapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/user"
)

// UserInfoFormat ...
type UserInfoFormat struct {
	user *userapi.UserGetResponse
}

// UsersInfoFormat ...
type UsersInfoFormat struct {
	users []*userapi.UserGetResponse
}

// NewUserInfoFormat ...
func NewUserInfoFormat(user *userapi.UserGetResponse) *UserInfoFormat {
	return &UserInfoFormat{
		user: user,
	}
}

// NewUsersInfoFormat ...
func NewUsersInfoFormat(users []*userapi.UserGetResponse) *UsersInfoFormat {
	return &UsersInfoFormat{
		users: users,
	}
}

// ToText ...
func (f *UserInfoFormat) ToText() (string, error) {
	res := fmt.Sprintf("ID:           %s\n", f.user.ID)
	res += fmt.Sprintf("Name:         %s\n", f.user.Name)
	res += fmt.Sprintf("Created Time: %s\n", f.user.CreatedAt)
	res += fmt.Sprintf("System Roles: %v\n", f.user.SystemRoles)
	res += fmt.Sprintf("Custom Roles: %v\n", f.user.CustomRoles)
	return res, nil
}

// ToJSON ...
func (f *UserInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.user)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToText ...
func (f *UsersInfoFormat) ToText() (string, error) {
	res := ""
	for i, prj := range f.users {
		format := NewUserInfoFormat(prj)
		msg, err := format.ToText()
		if err != nil {
			return "", err
		}
		res += msg
		if i < len(f.users)-1 {
			res += "\n---\n"
		}
	}
	return res, nil
}

// ToJSON ...
func (f *UsersInfoFormat) ToJSON() (string, error) {
	bytes, err := json.Marshal(f.users)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
