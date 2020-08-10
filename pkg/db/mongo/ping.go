package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
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
func (p *PingHandler) Ping() *errors.Error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := p.dbClient.Ping(ctx, nil); err != nil {
		return errors.New("DB failed", "DB Ping failed: %v", err)
	}
	return nil
}
