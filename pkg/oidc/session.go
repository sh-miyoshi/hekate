package oidc

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
)

// StartLoginSession ...
func StartLoginSession(projectName string, req *AuthRequest) (string, error) {
	resMode := req.ResponseMode
	if resMode == "" {
		resMode = "fragment"
		if len(req.ResponseType) == 1 {
			if req.ResponseType[0] == "code" || req.ResponseType[0] == "none" {
				resMode = "query"
			}
		}
	}

	s := &model.AuthCodeSession{
		SessionID:    uuid.New().String(),
		ExpiresIn:    time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		ClientID:     req.ClientID,
		RedirectURI:  req.RedirectURI,
		Nonce:        req.Nonce,
		ProjectName:  projectName,
		MaxAge:       req.MaxAge,
		ResponseMode: resMode,
		ResponseType: req.ResponseType,
		Prompt:       req.Prompt,
	}
	// *) userID, code will be set in after

	if err := db.GetInst().AuthCodeSessionAdd(s); err != nil {
		return "", errors.Wrap(err, "add auth code session failed")
	}
	return s.SessionID, nil
}

// VerifySession ...
func VerifySession(sessionID string) (*model.AuthCodeSession, error) {
	s, err := db.GetInst().AuthCodeSessionGet(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "user login session get failed")
	}

	// verify session
	now := time.Now().Unix()
	if now > s.ExpiresIn.Unix() {
		return nil, ErrSessionExpired
	}

	return s, nil
}

// TODO(DeleteSession)
