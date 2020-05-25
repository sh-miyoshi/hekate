package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
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
func NewClient(connStr string) (*mongo.Client, error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, errors.Wrap(err, "Mongo NewClient failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := cli.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "Mongo Connect failed")
	}

	if err := cli.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "Mongo Ping failed")
	}

	return cli, nil
}
