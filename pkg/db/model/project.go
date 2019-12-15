package model

import (
	"errors"
	"time"
)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  int
	RefreshTokenLifeSpan int
}

// ProjectInfo ...
type ProjectInfo struct {
	Name        string
	CreatedAt   time.Time
	TokenConfig *TokenConfig
}

var (
	// ErrProjectAlreadyExists ...
	ErrProjectAlreadyExists = errors.New("Project Already Exists")

	// ErrNoSuchProject ...
	ErrNoSuchProject = errors.New("No such project")
)
