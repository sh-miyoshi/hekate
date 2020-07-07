package memory

import "github.com/sh-miyoshi/hekate/pkg/errors"

// PingHandler implement db.PingHandler
type PingHandler struct {
}

// NewPingHandler ...
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping ...
func (p *PingHandler) Ping() *errors.Error {
	return nil
}
