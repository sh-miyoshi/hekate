package model

import (
	"errors"
	"time"
	"regexp"
)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  uint
	RefreshTokenLifeSpan uint
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

// Validate ...
func (p *ProjectInfo) Validate() error {
	prjNameRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-]{2,31}$`)
	if !prjNameRegExp.MatchString(p.Name) {
		return errors.New("Invalid Project Name format")
	}
	return nil
}