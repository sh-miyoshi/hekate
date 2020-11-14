package model

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Device ...
type Device struct {
	DeviceCode     string
	UserCode       string
	ProjectName    string
	ExpiresIn      int64
	CreatedAt      time.Time
	LoginSessionID string
}

// DeviceHandler ...
type DeviceHandler interface {
	Add(projectName string, ent *Device) *errors.Error
	DeleteAll(projectName string) *errors.Error
	Cleanup(now time.Time) *errors.Error

	// TODO add other methods
}

var (
	// ErrDeviceValidateFailed ...
	ErrDeviceValidateFailed = errors.New("Device validation failed", "Device validation failed")
)

// Validate ...
func (d *Device) Validate() *errors.Error {
	if d.DeviceCode == "" {
		return errors.Append(ErrDeviceValidateFailed, "DeviceCode is empty")
	}

	if d.UserCode == "" {
		return errors.Append(ErrDeviceValidateFailed, "UserCode is empty")
	}

	if !ValidateProjectName(d.ProjectName) {
		return errors.Append(ErrDeviceValidateFailed, "Invalid Project Name format")
	}

	if d.ExpiresIn <= 0 {
		return errors.Append(ErrDeviceValidateFailed, "expires time must be positive number, but got %d", d.ExpiresIn)
	}

	return nil
}
