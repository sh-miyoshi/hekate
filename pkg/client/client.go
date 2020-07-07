package client

import (
	"github.com/sh-miyoshi/hekate/pkg/db"
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
