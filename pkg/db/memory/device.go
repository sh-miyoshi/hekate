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

// GetList ...
func (h *DeviceHandler) GetList(projectName string, filter *model.DeviceFilter) ([]*model.Device, *errors.Error) {
	res := []*model.Device{}

	for _, role := range h.devices {
		if role.ProjectName == projectName {
			res = append(res, role)
		}
	}

	if filter != nil {
		res = matchFilterDeviceList(res, projectName, filter)
	}

	return res, nil
}

// Delete ...
func (h *DeviceHandler) Delete(projectName string, deviceCode string) *errors.Error {
	newList := []*model.Device{}
	found := false
	for _, d := range h.devices {
		if d.ProjectName == projectName && d.DeviceCode == deviceCode {
			found = true
		} else {
			newList = append(newList, d)
		}
	}

	if found {
		h.devices = newList
		return nil
	}
	return errors.New("Internal Error", "No such device %s", deviceCode)
}

// matchFilterDeviceList returns a list which matches the filter rules
func matchFilterDeviceList(data []*model.Device, projectName string, filter *model.DeviceFilter) []*model.Device {
	if filter == nil {
		return data
	}
	res := []*model.Device{}

	for _, d := range data {
		if projectName == d.ProjectName {
			if filter.DeviceCode != "" && d.DeviceCode != filter.DeviceCode {
				// missmatch device code
				continue
			}

			if filter.UserCode != "" && d.UserCode != filter.UserCode {
				// missmatch device code
				continue
			}
		}

		res = append(res, d)
	}

	return res
}
