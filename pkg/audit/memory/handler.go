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
func (h *Handler) Save(projectName string, tm time.Time, resType, method, path, message string) *errors.Error {
	h.data = append(h.data, model.Audit{
		ProjectName:  projectName,
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
func (h *Handler) Get(projectName string, fromDate, toDate time.Time) ([]model.Audit, *errors.Error) {
	res := []model.Audit{}

	// if we want to get logs whose date are from "2019-09-19",
	// we have to pass "2019-09-20 00:00:00.000" to mongodb.
	toDate = toDate.AddDate(0, 0, 1)

	for _, d := range h.data {
		if d.ProjectName == projectName {
			if fromDate.After(d.Time) && toDate.Before(d.Time) {
				res = append(res, d)
			}
		}
	}

	return res, nil
}
