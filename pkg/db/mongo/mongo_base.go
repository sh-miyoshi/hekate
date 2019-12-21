package mongo

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	databaseName          = "jwtserver"
	projectCollectionName = "project"

	timeoutSecond = 5
)

// NewClient ...
func NewClient(connStr string) (*mongo.Client, error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, errors.Cause(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := cli.Connect(ctx); err != nil {
		return nil, errors.Cause(err)
	}

	if err := cli.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "Failed to connect to mongodb")
	}

	return cli, nil
}
