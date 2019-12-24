package mongo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	dbClient *mongo.Client
}

// NewProjectHandler ...
func NewProjectHandler(dbClient *mongo.Client) (*ProjectInfoHandler, error) {
	res := &ProjectInfoHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name
	mod := mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(projectCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	v := &projectInfo{
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		TokenConfig: &tokenConfig{
			AccessTokenLifeSpan:  ent.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: ent.TokenConfig.RefreshTokenLifeSpan,
		},
	}

	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert project to mongodb")
	}

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete project from mongodb")
	}
	return nil
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]string, error) {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		return []string{}, errors.Wrap(err, "Failed to get project list from mongodb")
	}

	projects := []projectInfo{}
	if err := cursor.All(ctx, &projects); err != nil {
		return []string{}, errors.Wrap(err, "Failed to get project list from mongodb")
	}

	res := []string{}
	for _, prj := range projects {
		res = append(res, prj.Name)
	}

	return res, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.ProjectInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Cause(model.ErrNoSuchProject)
		}
		return nil, errors.Wrap(err, "Failed to get project from mongodb")
	}
	logger.Debug("Get project %s data: %v", name, res)

	return res, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: ent.Name},
	}

	v := &projectInfo{
		Name:      ent.Name,
		CreatedAt: ent.CreatedAt,
		TokenConfig: &tokenConfig{
			AccessTokenLifeSpan:  ent.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: ent.TokenConfig.RefreshTokenLifeSpan,
		},
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to update project in mongodb")
	}

	return nil
}
