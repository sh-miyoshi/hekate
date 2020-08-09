package model

import "time"

// Audit ...
type Audit struct {
	Time         time.Time
	ResourceType string
	Method       string
	Path         string
	Body         string
	// TODO(userID, clientID)
}
