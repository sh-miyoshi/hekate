package user

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

var (
	// ErrAuthFailed ...
	ErrAuthFailed = errors.New("Authentication failed", "Authentication failed")
	// ErrUserLocked ...
	ErrUserLocked = errors.New("User locked", "User locked")
)

func isLocked(state model.LockState, setting model.UserLock) bool {
	if !setting.Enabled {
		return false
	}

	if state.Locked {
		now := time.Now()
		last := state.VerifyFailedTimes[len(state.VerifyFailedTimes)-1]
		// If it's not yet time to unlock, return locked
		if now.Before(last.Add(setting.FailureResetTime)) {
			return true
		}
	}

	return false
}

func inclementFailedNum(state *model.LockState, setting model.UserLock) {
	if !setting.Enabled {
		return
	}

	now := time.Now()
	tmp := []time.Time{}

	// remove too old data
	old := now.Add(-setting.LockDuration)
	for _, t := range state.VerifyFailedTimes {
		if t.After(old) {
			tmp = append(tmp, t)
		}
	}

	// add now date
	state.VerifyFailedTimes = append(tmp, now)

	// update locked
	if len(state.VerifyFailedTimes) >= int(setting.MaxLoginFailure) {
		state.Locked = true
	}
}

// Verify ...
func Verify(projectName string, name string, password string) (*model.UserInfo, *errors.Error) {
	prj, err := db.GetInst().ProjectGet(projectName)
	if err != nil {
		return nil, err
	}
	users, err := db.GetInst().UserGetList(projectName, &model.UserFilter{Name: name})
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, ErrAuthFailed
	}
	user := users[0]

	if isLocked(user.LockState, prj.UserLock) {
		return nil, ErrUserLocked
	}

	hash := util.CreateHash(password)
	if user.PasswordHash != hash {
		// update lock state
		inclementFailedNum(&user.LockState, prj.UserLock)
		logger.Debug("user lock state: %v", user.LockState)
		if err := db.GetInst().UserUpdate(projectName, user); err != nil {
			return nil, err
		}
		return nil, ErrAuthFailed
	}

	// clear lock state
	if prj.UserLock.Enabled {
		logger.Debug("successfully user verify, so clear lock state")
		user.LockState = model.LockState{}
		if err := db.GetInst().UserUpdate(projectName, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
