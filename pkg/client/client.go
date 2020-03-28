package client

import (
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db"
)

var (
	// ErrNoRedirectURL ...
	ErrNoRedirectURL = errors.New("no such redirect url")
)

// CheckRedirectURL ...
func CheckRedirectURL(clientID, redirectURL string) error {
	// Check Redirect URL
	cli, err := db.GetInst().ClientGet(clientID)
	if err != nil {
		return err
	}

	for _, u := range cli.AllowedCallbackURLs {
		if u == redirectURL {
			return nil // found
		}
	}

	return ErrNoRedirectURL
}
