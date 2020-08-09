package memory

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Handler ...
type Handler struct {
	data []model.Audit
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

// Ping ...
func (h *Handler) Ping() *errors.Error {
	return nil
}

// Save ...
func (h *Handler) Save(tm time.Time, resType, method, path, message string) *errors.Error {
	h.data = append(h.data, model.Audit{
		Time:         tm,
		ResourceType: resType,
		Method:       method,
		Path:         path,
		IsSuccess:    message == "",
		Message:      message,
	})

	return nil
}

// Get ...
func (h *Handler) Get(fromDate, toDate time.Time) ([]model.Audit, *errors.Error) {
	res := []model.Audit{}

	for _, d := range h.data {
		if fromDate.After(d.Time) && toDate.Before(d.Time) {
			res = append(res, d)
		}
	}

	return res, nil
}
