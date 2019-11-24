package model

import (
	"errors"
)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  int32
	RefreshTokenLifeSpan int32
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
)
