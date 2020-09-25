package db

import (
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

var (
	// gcInterval is a interval time of running database garbage collector.
	// default: 1 hour
	gcInterval = 1 * time.Hour
)

// InitGC ...
func InitGC(intervalSec uint64) {
	gcInterval = time.Duration(intervalSec) * time.Second
}

// RunGC ...
func RunGC() {
	for {
		time.Sleep(gcInterval)
		logger.Debug("Delete expired sessions")
		if err := GetInst().DeleteExpiredSessions(); err != nil {
			errors.Print(errors.Append(err, "Failed to delete expired sessions"))
		}
	}
}
