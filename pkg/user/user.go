package user

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
)

var (
	// ErrAuthFailed ...
	ErrAuthFailed = errors.New("Authentication failed")
)

// Verify ...
func Verify(projectName string, name string, password string) (*model.UserInfo, error) {
	user, err := db.GetInst().UserGetByName(projectName, name)
	if err != nil {
		if errors.Cause(err) == model.ErrNoSuchUser {
			return nil, ErrAuthFailed
		}
		return nil, err
	}

	hash := util.CreateHash(password)
	if user.PasswordHash != hash {
		return nil, ErrAuthFailed
	}

	return user, nil
}
