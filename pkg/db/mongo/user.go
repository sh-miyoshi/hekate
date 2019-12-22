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
	v := &userInfo{
		ID:           ent.ID,
		ProjectName:  ent.ProjectName,
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PasswordHash: ent.PasswordHash,
		Roles:        ent.Roles,
		Sessions:     []session{},
	}

	for _, s := range ent.Sessions {
		v.Sessions = append(v.Sessions, session{
			SessionID: s.SessionID,
			CreatedAt: s.CreatedAt,
			ExpiresIn: s.ExpiresIn,
			FromIP:    s.FromIP,
		})
	}

	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert user to mongodb")
	}

	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete user from mongodb")
	}
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to get user list from mongodb")
	}

	users := []userInfo{}
	if err := cursor.All(ctx, &users); err != nil {
		return []string{}, errors.Wrap(err, "Failed to get user list from mongodb")
	}

	res := []string{}
	for _, user := range users {
		res = append(res, user.ID)
	}

	return res, nil
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
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: ent.ProjectName},
		{Key: "id", Value: ent.ID},
	}

	v := &userInfo{
		ID:           ent.ID,
		ProjectName:  ent.ProjectName,
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PasswordHash: ent.PasswordHash,
		Roles:        ent.Roles,
		Sessions:     []session{},
	}

	for _, s := range ent.Sessions {
		v.Sessions = append(v.Sessions, session{
			SessionID: s.SessionID,
			CreatedAt: s.CreatedAt,
			ExpiresIn: s.ExpiresIn,
			FromIP:    s.FromIP,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, v); err != nil {
		return errors.Wrap(err, "Failed to update project in mongodb")
	}

	return nil
}

// GetByName ...
func (h *UserInfoHandler) GetByName(projectName string, userName string) (*model.UserInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "name", Value: userName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.UserInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Cause(model.ErrNoSuchUser)
		}
		return nil, errors.Wrap(err, "Failed to get project by name from mongodb")
	}

	return res, nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete user from mongodb")
	}
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
