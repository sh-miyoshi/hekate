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

// CustomRoleHandler implement db.CustomRoleHandler
type CustomRoleHandler struct {
	session  mongo.Session
	dbClient *mongo.Client
}

// NewCustomRoleHandler ...
func NewCustomRoleHandler(dbClient *mongo.Client) (*CustomRoleHandler, error) {
	res := &CustomRoleHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name
	mod := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(roleCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *CustomRoleHandler) Add(ent *model.CustomRole) error {
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
		return errors.Wrap(err, "Failed to insert role to mongodb")
	}

	return nil
}

// Delete ...
func (h *CustomRoleHandler) Delete(roleID string) error {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)
	filter := bson.D{
		{Key: "id", Value: roleID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete role from mongodb")
	}
	return nil
}

// GetList ...
func (h *CustomRoleHandler) GetList(projectName string) ([]string, error) {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)

	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to get role list from mongodb")
	}

	roles := []customRole{}
	if err := cursor.All(ctx, &roles); err != nil {
		return []string{}, errors.Wrap(err, "Failed to get role list from mongodb")
	}

	res := []string{}
	for _, role := range roles {
		res = append(res, role.ID)
	}

	return res, nil
}

// Get ...
func (h *CustomRoleHandler) Get(roleID string) (*model.CustomRole, error) {
	col := h.dbClient.Database(databaseName).Collection(roleCollectionName)
	filter := bson.D{
		{Key: "id", Value: roleID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.CustomRole{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchCustomRole
		}
		return nil, errors.Wrap(err, "Failed to get role from mongodb")
	}

	return res, nil
}

// Update ...
func (h *CustomRoleHandler) Update(ent *model.CustomRole) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
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
		return errors.Wrap(err, "Failed to update client in mongodb")
	}

	return nil
}

// DeleteAll ...
func (h *CustomRoleHandler) DeleteAll(projectName string) error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete client from mongodb")
	}
	return nil
}

// BeginTx ...
func (h *CustomRoleHandler) BeginTx() error {
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
func (h *CustomRoleHandler) CommitTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.CommitTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}

// AbortTx ...
func (h *CustomRoleHandler) AbortTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.AbortTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}
