package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeviceHandler implement db.DeviceHandler
type DeviceHandler struct {
	dbClient *mongo.Client
}

// NewDeviceHandler ...
func NewDeviceHandler(dbClient *mongo.Client) (*DeviceHandler, *errors.Error) {
	res := &DeviceHandler{
		dbClient: dbClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(deviceCollectionName)
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
		logger.Info("Create index for device")
		// Create Index to Project Name and Device Code
		mod := mongo.IndexModel{
			Keys: bson.M{
				"project_name": 1, // index in ascending order
				"device_code":  1, // index in ascending order
			},
		}
		if _, err := iv.CreateOne(ctx, mod); err != nil {
			return nil, errors.New("DB failed", "Failed to create index: %v", err)
		}
	}

	return res, nil
}

// Add ...
func (h *DeviceHandler) Add(projectName string, ent *model.Device) *errors.Error {
	v := &device{
		DeviceCode:     ent.DeviceCode,
		UserCode:       ent.UserCode,
		ProjectName:    ent.ProjectName,
		ExpiresIn:      ent.ExpiresIn,
		CreatedAt:      ent.CreatedAt,
		LoginSessionID: ent.LoginSessionID,
	}

	col := h.dbClient.Database(databaseName).Collection(deviceCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("DB failed", "Failed to insert device to mongodb: %v", err)
	}

	return nil
}

// DeleteAll ...
func (h *DeviceHandler) DeleteAll(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(deviceCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete authcode session from mongodb: %v", err)
	}
	return nil
}

// Cleanup ...
func (h *DeviceHandler) Cleanup(now time.Time) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(deviceCollectionName)
	filter := bson.D{
		{Key: "expires_in", Value: bson.D{{Key: "$lt", Value: now.Unix()}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete expired device from mongodb: %v", err)
	}

	return nil
}
