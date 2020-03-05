package user

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

var (
	// ErrAuthFailed ...
	ErrAuthFailed = errors.New("Authentication failed")
)

// Verify ...
func Verify(projectName string, name string, password string) (*model.UserInfo, error) {
	users, err := db.GetInst().UserGetList(projectName, &model.UserFilter{Name: name})
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, ErrAuthFailed
	}
	user := users[0]

	hash := util.CreateHash(password)
	if user.PasswordHash != hash {
		return nil, ErrAuthFailed
	}

	return user, nil
}
