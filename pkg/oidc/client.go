package oidc

import (
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/stretchr/stew/slice"
)

var (
	// ErrNoRedirectURL ...
	ErrNoRedirectURL = errors.New("No such redirect url", "No such redirect url")
)

// CheckRedirectURL ...
func CheckRedirectURL(projectName, clientID, redirectURL string) *errors.Error {
	// Check Redirect URL
	cli, err := db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		return err
	}

	if ok := slice.Contains(cli.AllowedCallbackURLs, redirectURL); !ok {
		return ErrNoRedirectURL
	}

	return nil
}

// ClientAuth authenticates client with id and secret
func ClientAuth(projectName string, clientID string, clientSecret string) *errors.Error {
	client, err := db.GetInst().ClientGet(projectName, clientID)
	if err != nil {
		if errors.Contains(err, model.ErrNoSuchClient) || errors.Contains(err, model.ErrClientValidateFailed) {
			return errors.Append(errors.ErrInvalidClient, err.Error())
		}
		return errors.Append(err, "Failed to get client")
	}

	if client.AccessType != "public" {
		if client.Secret != clientSecret {
			return errors.Append(errors.ErrInvalidClient, "client auth failed")
		}
	}

	return nil
}
