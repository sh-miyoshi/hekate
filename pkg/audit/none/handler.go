package none

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Handler ...
type Handler struct{}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{}
}

// Ping ...
func (h *Handler) Ping() *errors.Error {
	return nil
}

// Save ...
func (h *Handler) Save(projectName string, tm time.Time, resType, method, path, message string) *errors.Error {
	return nil
}

// Get ...
func (h *Handler) Get(projectName string, fromDate, toDate time.Time, offset uint) ([]model.Audit, *errors.Error) {
	return []model.Audit{}, nil
}
