package memory

// PingHandler implement db.PingHandler
type PingHandler struct {
}

// NewPingHandler ...
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

// Ping ...
func (p *PingHandler) Ping() error {
	return nil
}
