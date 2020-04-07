package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	dbClient *mongo.Client
}

// NewClientHandler ...
func NewClientHandler(dbClient *mongo.Client) *ClientInfoHandler {
	res := &ClientInfoHandler{
		dbClient: dbClient,
	}

	// Client has no index

	return res
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
func (h *ClientInfoHandler) Delete(projectName, clientID string) error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
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
func (h *ClientInfoHandler) GetList(projectName string) ([]*model.ClientInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)

	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get client list from mongodb")
	}

	clients := []clientInfo{}
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, errors.Wrap(err, "Failed to get client list from mongodb")
	}

	res := []*model.ClientInfo{}
	for _, client := range clients {
		res = append(res, &model.ClientInfo{
			ID:                  client.ID,
			ProjectName:         client.ProjectName,
			Secret:              client.Secret,
			AccessType:          client.AccessType,
			CreatedAt:           client.CreatedAt,
			AllowedCallbackURLs: client.AllowedCallbackURLs,
		})
	}

	return res, nil
}

// Get ...
func (h *ClientInfoHandler) Get(projectName, clientID string) (*model.ClientInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "id", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &clientInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchClient
		}
		return nil, errors.Wrap(err, "Failed to get client from mongodb")
	}

	return &model.ClientInfo{
		ID:                  res.ID,
		ProjectName:         res.ProjectName,
		Secret:              res.Secret,
		AccessType:          res.AccessType,
		CreatedAt:           res.CreatedAt,
		AllowedCallbackURLs: res.AllowedCallbackURLs,
	}, nil
}

// Update ...
func (h *ClientInfoHandler) Update(ent *model.ClientInfo) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: ent.ProjectName},
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
