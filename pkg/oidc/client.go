package oidc

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
)

// ClientAuth ...
func ClientAuth(clientID string, clientSecret string) error {
	client, err := db.GetInst().ClientGet(clientID)
	if err != nil {
		e := errors.Cause(err)
		if e == model.ErrNoSuchClient || e == model.ErrClientValidateFailed {
			return ErrInvalidClient
		}
		return errors.Wrap(err, "Failed to get client")
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return ErrInvalidClient
		}
	}

	return nil
}
