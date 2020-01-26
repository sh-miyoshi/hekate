package oidc

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db"
)

// ClientAuth ...
func ClientAuth(clientID string, clientSecret string) error {
	client, err := db.GetInst().ClientGet(clientID)
	if err != nil {
		return err
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return ErrClientAuthFailed
		}
	}

	return nil
}
