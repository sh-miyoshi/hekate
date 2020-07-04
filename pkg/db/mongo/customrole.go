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

// CustomRoleHandler implement db.CustomRoleHandler
type CustomRoleHandler struct {
	dbClient *mongo.Client
}

// NewCustomRoleHandler ...
func NewCustomRoleHandler(dbClient *mongo.Client) (*CustomRoleHandler, *errors.Error) {
	res := &CustomRoleHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name and Custom Role ID
	mod := mongo.IndexModel{
		Keys: bson.M{
			"projectName": 1, // index in ascending order
			"id":          1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(roleCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)
	if err != nil {
		return nil, errors.New("Failed to create index: %v", err)
	}

	return res, nil
}

// Add ...
func (h *CustomRoleHandler) Add(projectName string, ent *model.CustomRole) *errors.Error {
	v := &customRole{
		ID:          ent.ID,
		ProjectName: ent.ProjectName,
		CreatedAt:   ent.CreatedAt,
		Name:        ent.Name,
	}

	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("Failed to insert role to mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *CustomRoleHandler) Delete(projectName string, roleID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: roleID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("Failed to delete role from mongodb: %v", err)
	}
	return nil
}

// GetList ...
func (h *CustomRoleHandler) GetList(projectName string, filter *model.CustomRoleFilter) ([]*model.CustomRole, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)

	f := bson.D{
		{Key: "projectName", Value: projectName},
	}

	if filter != nil {
		if filter.Name != "" {
			f = append(f, bson.E{Key: "name", Value: filter.Name})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
	if err != nil {
		return nil, errors.New("Failed to get custom role list from mongodb: %v", err)
	}

	roles := []customRole{}
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, errors.New("Failed to get custom role list from mongodb: %v", err)
	}

	res := []*model.CustomRole{}
	for _, role := range roles {
		res = append(res, &model.CustomRole{
			ID:          role.ID,
			ProjectName: role.ProjectName,
			CreatedAt:   role.CreatedAt,
			Name:        role.Name,
		})
	}

	return res, nil
}

// Get ...
func (h *CustomRoleHandler) Get(projectName string, roleID string) (*model.CustomRole, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: roleID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &customRole{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchCustomRole
		}
		return nil, errors.New("Failed to get role from mongodb: %v", err)
	}

	return &model.CustomRole{
		ID:          res.ID,
		ProjectName: res.ProjectName,
		CreatedAt:   res.CreatedAt,
		Name:        res.Name,
	}, nil
}

// Update ...
func (h *CustomRoleHandler) Update(projectName string, ent *model.CustomRole) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: ent.ID},
	}

	v := &customRole{
		ID:          ent.ID,
		ProjectName: ent.ProjectName,
		CreatedAt:   ent.CreatedAt,
		Name:        ent.Name,
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.New("Failed to update client in mongodb: %v", err)
	}

	return nil
}

// DeleteAll ...
func (h *CustomRoleHandler) DeleteAll(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("Failed to delete client from mongodb: %v", err)
	}
	return nil
}
