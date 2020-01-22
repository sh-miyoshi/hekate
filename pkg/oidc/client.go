package oidc

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
)

var (
	// ErrAuthFailed ...
	ErrAuthFailed = errors.New("missing client secret")
)

// ClientAuth ...
func ClientAuth(clientID string, clientSecret string) error {
	client, err := db.GetInst().ClientGet(clientID)
	if err != nil {
		return err
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return ErrAuthFailed
		}
	}

	return nil
}
