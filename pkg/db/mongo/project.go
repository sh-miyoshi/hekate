package mongo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	dbClient *mongo.Client
}

// NewProjectHandler ...
func NewProjectHandler(dbClient *mongo.Client) *ProjectInfoHandler {
	res := &ProjectInfoHandler{
		dbClient: dbClient,
	}

	return res
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, ent)
	if err != nil {
		return errors.Wrap(err, "Failed to insert project to mongodb")
	}

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) error {
	return nil
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]string, error) {
	return []string{}, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, error) {
	return nil, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
	return nil
}