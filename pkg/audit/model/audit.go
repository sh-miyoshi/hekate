package model

import "time"

// Audit ...
type Audit struct {
	ProjectName  string
	Time         time.Time
	ResourceType string
	Method       string
	Path         string
	IsSuccess    bool
	Message      string
	// TODO(userID, clientID)
}
