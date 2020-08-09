package model

import "time"

// Audit ...
type Audit struct {
	Time         time.Time
	ResourceType string
	Method       string
	Path         string
	IsSuccess    bool
	Message      string
	// TODO(userID, clientID)
}
