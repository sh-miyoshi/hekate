package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
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

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(sessionCollectionName)
	iv := col.Indexes()
	var ires []bson.M
	cur, err := iv.List(ctx)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get index info: %v", err)
	}
	if err := cur.All(ctx, &ires); err != nil {
		return nil, errors.New("DB failed", "Failed to get index info: %v", err)
	}

	if len(ires) == 0 {
		logger.Info("Create index for session")
		// Create Index to Project Name and SessionID
		mod := mongo.IndexModel{
			Keys: bson.M{
				"project_name": 1, // index in ascending order
				"session_id":   1, // index in ascending order
			},
			Options: options.Index().SetUnique(true),
		}
		if _, err := iv.CreateOne(ctx, mod); err != nil {
			return nil, errors.New("DB failed", "Failed to create index: %v", err)
		}
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
func (h *SessionHandler) Delete(projectName string, filter *model.SessionFilter) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)
	f := bson.D{
		{Key: "project_name", Value: projectName},
	}
	if filter != nil {
		if filter.SessionID != "" {
			f = append(f, bson.E{Key: "session_id", Value: filter.SessionID})
		}
		if filter.UserID != "" {
			f = append(f, bson.E{Key: "user_id", Value: filter.UserID})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, f)
	if err != nil {
		return errors.New("DB failed", "Failed to delete session from mongodb: %v", err)
	}
	return nil
}

// DeleteAll ...
func (h *SessionHandler) DeleteAll(projectName string) *errors.Error {
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

// GetList ...
func (h *SessionHandler) GetList(projectName string, filter *model.SessionFilter) ([]*model.Session, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(sessionCollectionName)

	f := bson.D{
		{Key: "project_name", Value: projectName},
	}
	if filter != nil {
		if filter.SessionID != "" {
			f = append(f, bson.E{Key: "session_id", Value: filter.SessionID})
		}
		if filter.UserID != "" {
			f = append(f, bson.E{Key: "user_id", Value: filter.UserID})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
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
