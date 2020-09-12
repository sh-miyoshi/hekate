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

// ClientInfoHandler implement db.ClientInfoHandler
type ClientInfoHandler struct {
	dbClient *mongo.Client
}

// NewClientHandler ...
func NewClientHandler(dbClient *mongo.Client) (*ClientInfoHandler, *errors.Error) {
	res := &ClientInfoHandler{
		dbClient: dbClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(clientCollectionName)
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
		logger.Info("Create index for client")
		// Create Index to Project Name and Client ID
		mod := mongo.IndexModel{
			Keys: bson.M{
				"project_name": 1, // index in ascending order
				"id":           1, // index in ascending order
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
func (h *ClientInfoHandler) Add(projectName string, ent *model.ClientInfo) *errors.Error {
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
		return errors.New("DB failed", "Failed to insert client to mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *ClientInfoHandler) Delete(projectName, clientID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "id", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete client from mongodb: %v", err)
	}
	return nil
}

// GetList ...
func (h *ClientInfoHandler) GetList(projectName string, filter *model.ClientFilter) ([]*model.ClientInfo, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)

	f := bson.D{
		{Key: "project_name", Value: projectName},
	}

	if filter != nil {
		if filter.ID != "" {
			f = append(f, bson.E{Key: "id", Value: filter.ID})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get client list from mongodb: %v", err)
	}

	clients := []clientInfo{}
	if err := cursor.All(ctx, &clients); err != nil {
		return nil, errors.New("DB failed", "Failed to get client list from mongodb: %v", err)
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

// Update ...
func (h *ClientInfoHandler) Update(projectName string, ent *model.ClientInfo) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
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
		return errors.New("DB failed", "Failed to update client in mongodb: %v", err)
	}

	return nil
}

// DeleteAll ...
func (h *ClientInfoHandler) DeleteAll(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(clientCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete client from mongodb: %v", err)
	}
	return nil
}
