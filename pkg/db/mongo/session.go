package mongo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	session  mongo.Session
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

// RevokeAll ...
func (h *SessionHandler) RevokeAll(userID string) error {
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

// BeginTx ...
func (h *SessionHandler) BeginTx() error {
	var err error
	h.session, err = h.dbClient.StartSession()
	if err != nil {
		return err
	}
	err = h.session.StartTransaction()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
		defer cancel()
		h.session.EndSession(ctx)
		return err
	}
	return nil
}

// CommitTx ...
func (h *SessionHandler) CommitTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.CommitTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}

// AbortTx ...
func (h *SessionHandler) AbortTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.AbortTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}
