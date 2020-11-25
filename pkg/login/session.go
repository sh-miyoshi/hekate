package login

import (
	"time"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
)

// StartLoginSession ...
func StartLoginSession(projectName string, req *oidc.AuthRequest) (string, *errors.Error) {
	expires := time.Second * time.Duration(config.Get().LoginSessionExpiresIn)

	s := &model.LoginSession{
		SessionID:           uuid.New().String(),
		ExpiresDate:         time.Now().Add(expires),
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		Nonce:               req.Nonce,
		ProjectName:         projectName,
		ResponseMode:        req.ResponseMode,
		ResponseType:        req.ResponseType,
		Prompt:              req.Prompt,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
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
	now := time.Now()
	if now.After(s.ExpiresDate) {
		return nil, errors.ErrSessionExpired
	}

	return s, nil
}
