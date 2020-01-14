package oidc

import (
	"encoding/json"
	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/google/uuid"
	"time"
)

var (
	expiresTimeSec = uint64(10 * 60) // default: 10 minutes
)

type authCode struct {
	ExpiresIn   time.Time `json:"expires_in"`
	ClientID    string    `json:"client_id"`
	RedirectURL string    `json:"redirect_url"`
	CodeID      string    `json:"code_id"`
}

// InitAuthCodeConfig ...
func InitAuthCodeConfig(authCodeExpiresTimeSec uint64) {
	expiresTimeSec = authCodeExpiresTimeSec
}

// GenerateAuthCode ...
func GenerateAuthCode(clientID string, redirectURL string) string {
	code := authCode{
		CodeID:      uuid.New().String(),
		ClientID:    clientID,
		RedirectURL: redirectURL,
		ExpiresIn:   time.Now().Add(time.Second * time.Duration(expiresTimeSec)),
	}
	b, _ := json.Marshal(code)
	return base64url.Encode(b)
}
