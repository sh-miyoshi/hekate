package model

import (
	"errors"
)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  int
	RefreshTokenLifeSpan int
}

// ProjectInfo ...
type ProjectInfo struct {
	ID          string
	Name        string
	Enabled     bool
	CreatedAt   string
	TokenConfig *TokenConfig
}

var (
	// ErrProjectAlreadyExists ...
	ErrProjectAlreadyExists = errors.New("Project Already Exists")

	// ErrNoSuchProject ...
	ErrNoSuchProject = errors.New("No such project")
)
