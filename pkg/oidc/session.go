package oidc

import (
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// StartLoginSession ...
func StartLoginSession(projectName string, req *AuthRequest) (string, *errors.Error) {
	s := &model.LoginSession{
		SessionID:    uuid.New().String(),
		ExpiresIn:    time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		ClientID:     req.ClientID,
		RedirectURI:  req.RedirectURI,
		Nonce:        req.Nonce,
		ProjectName:  projectName,
		MaxAge:       req.MaxAge,
		ResponseMode: req.ResponseMode,
		ResponseType: req.ResponseType,
		Prompt:       req.Prompt,
	}
	// *) userID, code will be set in after

	if err := db.GetInst().LoginSessionAdd(projectName, s); err != nil {
		return "", errors.Append(err, "add login session failed")
	}
	return s.SessionID, nil
}

// VerifySession ...
func VerifySession(projectName string, sessionID string) (*model.LoginSession, *errors.Error) {
	s, err := db.GetInst().LoginSessionGet(projectName, sessionID)
	if err != nil {
		return nil, errors.Append(err, "user login session get failed")
	}

	// verify session
	now := time.Now().Unix()
	if now > s.ExpiresIn.Unix() {
		return nil, errors.ErrSessionExpired
	}

	return s, nil
}
