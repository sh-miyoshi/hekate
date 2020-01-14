package oidc

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
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
func GenerateAuthCode(clientID string, redirectURL string) (string, error) {
	code := &model.AuthCode{
		CodeID:      uuid.New().String(),
		ClientID:    clientID,
		RedirectURL: redirectURL,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
	}

	err := db.GetInst().NewAuthCode(code)

	return code.CodeID, err
}
