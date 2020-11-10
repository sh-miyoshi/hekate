package audit

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

// Query ...
type Query struct {
	FromDate time.Time
	ToDate   time.Time
	Offset   uint
}

// Handler ...
type Handler interface {
	Ping() *errors.Error
	Save(projectName string, tm time.Time, resType, method, path, message string) *errors.Error
	Get(projectName string, fromDate, toDate time.Time, offset uint) ([]model.Audit, *errors.Error)
}
