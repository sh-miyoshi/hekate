package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName                  = "hekate"
	projectCollectionName         = "project"
	userCollectionName            = "user"
	clientCollectionName          = "client"
	sessionCollectionName         = "session"
	authCodeCollectionName        = "code"
	roleCollectionName            = "customrole"
	authcodeSessionCollectionName = "authcodesession"
	roleInUserCollectionName      = "customroleinuser"

	timeoutSecond = 5
)

// NewClient ...
func NewClient(connStr string) (*mongo.Client, *errors.Error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, errors.New("", "Failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := cli.Connect(ctx); err != nil {
		return nil, errors.New("", "Failed to connect to mongo: %v", err)
	}

	if err := cli.Ping(ctx, nil); err != nil {
		return nil, errors.New("", "Failed to ping to mongo: %v", err)
	}

	return cli, nil
}
