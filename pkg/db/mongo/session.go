package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	dbClient *mongo.Client
}

// NewSessionHandler ...
func NewSessionHandler(dbClient *mongo.Client) (*SessionHandler, *errors.Error) {
	res := &SessionHandler{
		dbClient: dbClient,
	}

	// Create Index to ProjectName and SessionID
	mod := mongo.IndexModel{
		Keys: bson.M{
			"project_name": 1, // index in ascending order
			"session_id":   1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(sessionCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to create index: %v", err)
	}

	return res, nil
}

// Add ...
func (h *SessionHandler) Add(projectName string, s *model.Session) *errors.Error {
	v := &session{
		UserID:       s.UserID,
		ProjectName:  s.ProjectName,
		SessionID:    s.SessionID,
		CreatedAt:    s.CreatedAt,
		ExpiresIn:    s.ExpiresIn,
		FromIP:       s.FromIP,
		LastAuthTime: s.LastAuthTime,
	}

	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("DB failed", "Failed to insert session to mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *SessionHandler) Delete(projectName string, sessionID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "session_id", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete session from mongodb: %v", err)
	}
	return nil
}

// DeleteAll ...
func (h *SessionHandler) DeleteAll(projectName string, userID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "user_id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete session from mongodb: %v", err)
	}
	return nil
}

// DeleteAllInProject ...
func (h *SessionHandler) DeleteAllInProject(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete session from mongodb: %v", err)
	}
	return nil
}

// Get ...
func (h *SessionHandler) Get(projectName string, sessionID string) (*model.Session, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "session_id", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &session{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchSession
		}
		return nil, errors.New("DB failed", "Failed to get session from mongodb: %v", err)
	}

	return &model.Session{
		UserID:       res.UserID,
		ProjectName:  res.ProjectName,
		SessionID:    res.SessionID,
		CreatedAt:    res.CreatedAt,
		ExpiresIn:    res.ExpiresIn,
		FromIP:       res.FromIP,
		LastAuthTime: res.LastAuthTime,
	}, nil
}

// GetList ...
func (h *SessionHandler) GetList(projectName string, userID string) ([]*model.Session, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "user_id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get session list from mongodb: %v", err)
	}

	sessions := []session{}
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, errors.New("DB failed", "Failed to parse session list from mongodb: %v", err)
	}

	res := []*model.Session{}
	for _, s := range sessions {
		res = append(res, &model.Session{
			UserID:       s.UserID,
			ProjectName:  s.ProjectName,
			SessionID:    s.SessionID,
			CreatedAt:    s.CreatedAt,
			ExpiresIn:    s.ExpiresIn,
			FromIP:       s.FromIP,
			LastAuthTime: s.LastAuthTime,
		})
	}

	return res, nil
}
