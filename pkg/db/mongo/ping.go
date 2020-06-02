package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// PingHandler implement db.PingHandler
type PingHandler struct {
	dbClient *mongo.Client
}

// NewPingHandler ...
func NewPingHandler(dbClient *mongo.Client) *PingHandler {
	res := &PingHandler{
		dbClient: dbClient,
	}

	return res
}

// Ping ...
func (p *PingHandler) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	return p.dbClient.Ping(ctx, nil)
}
