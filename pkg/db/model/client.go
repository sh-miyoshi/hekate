package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"regexp"
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
	// Check ID
	if !(2 <= len(c.ID) && len(c.ID) < 128) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client ID format")
	}

	// Check Project Name
	prjNameRegExp := regexp.MustCompile(`^[a-z][a-z0-9\-]{2,31}$`)
	if !prjNameRegExp.MatchString(c.ProjectName) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Project Name format")
	}

	// Check Secret
	if !(8 <= len(c.ID) && len(c.ID) < 256) {
		return errors.Wrap(ErrClientValidateFailed, "Invalid Client Secret format")
	}

	// Check Access Type
	if c.AccessType != "public" && c.AccessType != "confidential" {
		return errors.Wrap(ErrClientValidateFailed, "Invalid access type")
	}

	for _, u := range c.AllowedCallbackURLs {
		if ok := govalidator.IsRequestURL(u); !ok {
			return errors.Wrap(ErrClientValidateFailed, "Invalid callback URL")
		}
	}

	return nil
}
