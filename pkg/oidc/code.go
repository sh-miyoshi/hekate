package oidc

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"net/http"
	"time"
)

var (
	expiresTimeSec = uint64(10 * 60) // default: 10 minutes
)

// InitAuthCodeConfig ...
func InitAuthCodeConfig(authCodeExpiresTimeSec uint64) {
	expiresTimeSec = authCodeExpiresTimeSec
}

// GenerateAuthCode ...
func GenerateAuthCode(clientID string, redirectURL string, userID string) (string, error) {
	code := &model.AuthCode{
		CodeID:      uuid.New().String(),
		ClientID:    clientID,
		RedirectURL: redirectURL,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
		UserID:      userID,
	}

	err := db.GetInst().NewAuthCode(code)

	return code.CodeID, err
}

// ValidateAuthCode ...
func ValidateAuthCode(codeID string) (*model.AuthCode, int, string) {
	code, err := db.GetInst().GetAuthCode(codeID)
	if err != nil {
		if err == model.ErrNoSuchCode {
			// TODO(revoke all token in code.UserID) <- SHOULD
			logger.Info("No such code %s", codeID)
			return nil, http.StatusBadRequest, "No such code"
		}
		logger.Error("Failed to get auth code: %+v", err)
		return nil, http.StatusInternalServerError, "Internal Server Error"
	}
	logger.Debug("Code: %v", code)

	if time.Now().Unix() >= code.ExpiresIn.Unix() {
		logger.Info("code %s is expired at %v", codeID, code.ExpiresIn)
		return nil, http.StatusBadRequest, "Code expired"
	}

	return code, http.StatusOK, "ok"
}
