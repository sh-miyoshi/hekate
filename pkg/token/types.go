package token

import (
	"time"
)

// Request ...
type Request struct {
	ExpiredTime time.Duration
	Audience    string
}
