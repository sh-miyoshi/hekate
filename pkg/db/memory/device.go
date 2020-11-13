package memory

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// DeviceHandler implement db.DeviceHandler
type DeviceHandler struct {
	devices []*model.Device
}

// NewDeviceHandler ...
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// Add ...
func (h *DeviceHandler) Add(projectName string, ent *model.Device) *errors.Error {
	h.devices = append(h.devices, ent)
	return nil
}

// DeleteAll ...
func (h *DeviceHandler) DeleteAll(projectName string) *errors.Error {
	newList := []*model.Device{}
	for _, s := range h.devices {
		if s.ProjectName != projectName {
			newList = append(newList, s)
		}
	}

	h.devices = newList
	return nil
}

// Cleanup ...
func (h *DeviceHandler) Cleanup(now time.Time) *errors.Error {
	newList := []*model.Device{}
	for _, s := range h.devices {
		expire := s.CreatedAt.Add(time.Second * time.Duration(s.ExpiresIn))
		if now.Before(expire) {
			newList = append(newList, s)
		}
	}

	h.devices = newList
	return nil
}
