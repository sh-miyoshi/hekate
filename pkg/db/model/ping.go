package model

import "github.com/sh-miyoshi/hekate/pkg/errors"

// PingHandler ...
type PingHandler interface {
	Ping() *errors.Error
}
