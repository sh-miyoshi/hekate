package client

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/stretchr/stew/slice"
)

var (
	// ErrNoRedirectURL ...
	ErrNoRedirectURL = errors.New("no such redirect url")
)

// CheckRedirectURL ...
func CheckRedirectURL(projectName, clientID, redirectURL string) error {
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
