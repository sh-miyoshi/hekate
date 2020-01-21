package mongo

import (
	"context"
	"github.com/pkg/errors"
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
func (h *SessionHandler) New(s *model.Session) error {
	v := &session{
		UserID:    s.UserID,
		SessionID: s.SessionID,
		CreatedAt: s.CreatedAt,
		ExpiresIn: s.ExpiresIn,
		FromIP:    s.FromIP,
	}

	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert session to mongodb")
	}

	return nil
}

// Revoke ...
func (h *SessionHandler) Revoke(sessionID string) error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "sessionID", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete session from mongodb")
	}
	return nil
}

// Get ...
func (h *SessionHandler) Get(sessionID string) (*model.Session, error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "sessionID", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.Session{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Cause(model.ErrNoSuchSession)
		}
		return nil, errors.Wrap(err, "Failed to get session from mongodb")
	}

	return res, nil
}

// GetList ...
func (h *SessionHandler) GetList(userID string) ([]string, error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	filter := bson.D{
		{Key: "userID", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to get session list from mongodb")
	}

	sessions := []session{}
	if err := cursor.All(ctx, &sessions); err != nil {
		return []string{}, errors.Wrap(err, "Failed to get session list from mongodb")
	}

	res := []string{}
	for _, s := range sessions {
		res = append(res, s.SessionID)
	}

	return res, nil
}
