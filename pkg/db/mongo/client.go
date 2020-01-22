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

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	session  mongo.Session
	dbClient *mongo.Client
}

// NewClientHandler ...
func NewClientHandler(dbClient *mongo.Client) (*ClientInfoHandler, error) {
	res := &ClientInfoHandler{
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

	col := res.dbClient.Database(databaseName).Collection(clientCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *ClientInfoHandler) Add(ent *model.ClientInfo) error {
	v := &clientInfo{
		ID:                  ent.ID,
		ProjectName:         ent.ProjectName,
		Secret:              ent.Secret,
		AccessType:          ent.AccessType,
		CreatedAt:           ent.CreatedAt,
		AllowedCallbackURLs: ent.AllowedCallbackURLs,
	}

	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert client to mongodb")
	}

	return nil
}

// Delete ...
func (h *ClientInfoHandler) Delete(clientID string) error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "id", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete client from mongodb")
	}
	return nil
}

// GetList ...
func (h *ClientInfoHandler) GetList(projectName string) ([]string, error) {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)

	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to get client list from mongodb")
	}

	clients := []clientInfo{}
	if err := cursor.All(ctx, &clients); err != nil {
		return []string{}, errors.Wrap(err, "Failed to get client list from mongodb")
	}

	res := []string{}
	for _, client := range clients {
		res = append(res, client.ID)
	}

	return res, nil
}

// Get ...
func (h *ClientInfoHandler) Get(clientID string) (*model.ClientInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "id", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.ClientInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchClient
		}
		return nil, errors.Wrap(err, "Failed to get client from mongodb")
	}

	return res, nil
}

// Update ...
func (h *ClientInfoHandler) Update(ent *model.ClientInfo) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "id", Value: ent.ID},
	}

	v := &clientInfo{
		ID:                  ent.ID,
		ProjectName:         ent.ProjectName,
		Secret:              ent.Secret,
		AccessType:          ent.AccessType,
		CreatedAt:           ent.CreatedAt,
		AllowedCallbackURLs: ent.AllowedCallbackURLs,
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
func (h *ClientInfoHandler) DeleteAll(projectName string) error {
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
func (h *ClientInfoHandler) BeginTx() error {
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
func (h *ClientInfoHandler) CommitTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.CommitTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}

// AbortTx ...
func (h *ClientInfoHandler) AbortTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.AbortTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}
