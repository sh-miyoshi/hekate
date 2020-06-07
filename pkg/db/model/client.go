package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
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
	Add(projectName string, ent *ClientInfo) error
	Delete(projectName, clientID string) error
	GetList(projectName string) ([]*ClientInfo, error)
	Get(projectName, clientID string) (*ClientInfo, error)
	Update(projectName string, ent *ClientInfo) error
	DeleteAll(projectName string) error
}

// Validate ...
func (c *ClientInfo) Validate() error {
	if !ValidateClientID(c.ID) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client ID format")
	}

	if !ValidateProjectName(c.ProjectName) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Project Name format")
	}

	if !ValidateClientAccessType(c.AccessType) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid access type")
	}

	if !ValidateClientSecret(c.Secret, c.AccessType) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client Secret format")
	}

	for _, u := range c.AllowedCallbackURLs {
		if !govalidator.IsRequestURL(u) {
			return errors.Wrap(ErrClientValidateFailed, "Invalid callback URL")
		}
	}

	return nil
}
