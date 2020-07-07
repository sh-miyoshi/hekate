package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/sh-miyoshi/hekate/pkg/errors"
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
	ErrClientAlreadyExists = errors.New("Client already exists", "Client already exists")

	// ErrNoSuchClient ...
	ErrNoSuchClient = errors.New("No such client", "No such client")

	// ErrClientValidateFailed ...
	ErrClientValidateFailed = errors.New("Client validation failed", "Client validation failed")
)

// ClientInfoHandler ...
type ClientInfoHandler interface {
	Add(projectName string, ent *ClientInfo) *errors.Error
	Delete(projectName, clientID string) *errors.Error
	GetList(projectName string) ([]*ClientInfo, *errors.Error)
	Get(projectName, clientID string) (*ClientInfo, *errors.Error)
	Update(projectName string, ent *ClientInfo) *errors.Error
	DeleteAll(projectName string) *errors.Error
}

// Validate ...
func (c *ClientInfo) Validate() *errors.Error {
	if !ValidateClientID(c.ID) {
		return errors.Append(ErrClientValidateFailed, "Invalid Client ID format")
	}

	if !ValidateProjectName(c.ProjectName) {
		return errors.Append(ErrClientValidateFailed, "Invalid Project Name format")
	}

	if !ValidateClientAccessType(c.AccessType) {
		return errors.Append(ErrClientValidateFailed, "Invalid access type")
	}

	if !ValidateClientSecret(c.Secret, c.AccessType) {
		return errors.Append(ErrClientValidateFailed, "Invalid Client Secret format")
	}

	for _, u := range c.AllowedCallbackURLs {
		if !govalidator.IsRequestURL(u) {
			return errors.Append(ErrClientValidateFailed, "Invalid callback URL")
		}
	}

	return nil
}
