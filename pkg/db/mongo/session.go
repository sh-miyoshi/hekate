package mongo

import (
	"context"
	//"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	dbClient *mongo.Client
}

// NewSessionHandler ...
func NewSessionHandler(dbClient *mongo.Client) (*SessionHandler, error) {
	res := &SessionHandler{
		dbClient: dbClient,
	}

	// Create Index to SessionID
	mod := mongo.IndexModel{
		Keys: bson.M{
			"sessionID": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(sessionCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// New ...
func (h *SessionHandler) New(userID string, sessionID string, expiresIn uint, fromIP string) error {
	return nil
}

// Revoke ...
func (h *SessionHandler) Revoke(sessionID string) error {
	return nil
}

// Get ...
func (h *SessionHandler) Get(sessionID string) (*model.Session, error) {
	return nil, nil
}

// GetList ...
func (h *SessionHandler) GetList(userID string) ([]string, error) {
	res := []string{}

	return res, nil
}
