package model

import (
	"errors"
	"time"
)

// ClientInfo ...
type ClientInfo struct {
	ID          string
	ProjectName string
	Secret      string
	AccessType  string
	CreatedAt   time.Time
}

var (
	// ErrClientAlreadyExists ...
	ErrClientAlreadyExists = errors.New("Client Already Exists")

	// ErrNoSuchClient ...
	ErrNoSuchClient = errors.New("No Such Client")
)

// ClientInfoHandler ...
type ClientInfoHandler interface {
	Add(ent *ClientInfo) error
	Delete(clientID string) error
	GetList(projectName string) ([]string, error)
	Get(clientID string) (*ClientInfo, error)
	Update(ent *ClientInfo) error
	DeleteAll(projectName string) error
}
