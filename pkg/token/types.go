package token

import (
	"time"
)

// Request ...
type Request struct {
	Issuer      string
	ExpiredTime time.Duration
	Audience    string
}
