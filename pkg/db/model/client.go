package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"time"
)

// ClientInfo ...
type ClientInfo struct {
	ID                  string
	ProjectName         string
	Secret              string
	AccessType          string
	CreatedAt           time.Time
	AllowedCallbackURLs []string
}

var (
	// ErrClientAlreadyExists ...
	ErrClientAlreadyExists = errors.New("Client Already Exists")

	// ErrNoSuchClient ...
	ErrNoSuchClient = errors.New("No Such Client")

	// ErrClientValidateFailed ...
	ErrClientValidateFailed = errors.New("Client validation failed")
)

// ClientInfoHandler ...
type ClientInfoHandler interface {
	Add(ent *ClientInfo) error
	Delete(clientID string) error
	GetList(projectName string) ([]string, error)
	Get(clientID string) (*ClientInfo, error)
	Update(ent *ClientInfo) error
	DeleteAll(projectName string) error

	// BeginTx method starts a transaction
	BeginTx() error

	// CommitTx method commits the transaction
	CommitTx() error

	// AbortTx method abort and rollback the transaction
	AbortTx() error
}

// Validate ...
func (c *ClientInfo) Validate() error {
	if !validateClientID(c.ID) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client ID format")
	}

	if !validateProjectName(c.ProjectName) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Project Name format")
	}

	if !validateClientSecret(c.Secret) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client Secret format")
	}

	if !validateClientAccessType(c.AccessType) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid access type")
	}

	for _, u := range c.AllowedCallbackURLs {
		if !govalidator.IsRequestURL(u) {
			return errors.Wrap(ErrClientValidateFailed, "Invalid callback URL")
		}
	}

	return nil
}
