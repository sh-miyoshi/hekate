package model

import (
	"time"
)

// OTP ...
type OTP struct {
	ID              string
	ProjectName     string
	PrivateKey      string
	InitExpiresDate time.Time
}
