package model

import (
	"errors"
)

// ProjectInfo ...
type ProjectInfo struct {
	ID   string
	Name string
}

var (
	// ErrProjectAlreadyExists ...
	ErrProjectAlreadyExists = errors.New("Project Already Exists")
)
