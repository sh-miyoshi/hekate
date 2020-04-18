package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Add ...
func (h *SessionHandler) Add(s *model.Session) error {
	v := &session{
		UserID:      s.UserID,
		ProjectName: s.ProjectName,
		SessionID:   s.SessionID,
		CreatedAt:   s.CreatedAt,
		ExpiresIn:   s.ExpiresIn,
		FromIP:      s.FromIP,
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

// Delete ...
func (h *SessionHandler) Delete(sessionID string) error {
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

// DeleteAll ...
func (h *SessionHandler) DeleteAll(userID string) error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "userID", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
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
			return nil, model.ErrNoSuchSession
		}
		return nil, errors.Wrap(err, "Failed to get session from mongodb")
	}

	return res, nil
}

// GetList ...
func (h *SessionHandler) GetList(userID string) ([]*model.Session, error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	filter := bson.D{
		{Key: "userID", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get session list from mongodb")
	}

	sessions := []session{}
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, errors.Wrap(err, "Failed to get session list from mongodb")
	}

	res := []*model.Session{}
	for _, s := range sessions {
		res = append(res, &model.Session{
			UserID:      s.UserID,
			ProjectName: s.ProjectName,
			SessionID:   s.SessionID,
			CreatedAt:   s.CreatedAt,
			ExpiresIn:   s.ExpiresIn,
			FromIP:      s.FromIP,
		})
	}

	return res, nil
}
