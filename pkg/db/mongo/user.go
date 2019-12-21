package mongo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	dbClient *mongo.Client
}

// NewUserHandler ...
func NewUserHandler(dbClient *mongo.Client) *UserInfoHandler {
	res := &UserInfoHandler{
		dbClient: dbClient,
	}
	return res
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) error {
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	return []string{}, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectName string, userID string) (*model.UserInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.UserInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Cause(model.ErrNoSuchUser)
		}
		return nil, errors.Wrap(err, "Failed to get project from mongodb")
	}

	return res, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	return nil
}

// GetByName ...
func (h *UserInfoHandler) GetByName(projectName string, userName string) (*model.UserInfo, error) {
	return nil, nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) error {
	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(projectName string, userID string, roleID string) error {
	return nil
}

// DeleteRole ...
func (h *UserInfoHandler) DeleteRole(projectName string, userID string, roleID string) error {
	return nil
}

// NewSession ...
func (h *UserInfoHandler) NewSession(projectName string, userID string, session model.Session) error {
	return nil
}

// RevokeSession ...
func (h *UserInfoHandler) RevokeSession(projectName string, userID string, sessionID string) error {
	return nil
}
