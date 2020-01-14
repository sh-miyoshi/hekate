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

// AuthCodeHandler implement db.AuthCodeHandler
type AuthCodeHandler struct {
	dbClient *mongo.Client
}

// NewAuthCodeHandler ...
func NewAuthCodeHandler(dbClient *mongo.Client) (*AuthCodeHandler, error) {
	res := &AuthCodeHandler{
		dbClient: dbClient,
	}

	// Create Index to AuthCodeID
	mod := mongo.IndexModel{
		Keys: bson.M{
			"codeID": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(authCodeCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// New ...
func (h *AuthCodeHandler) New(code *model.AuthCode) error {
	v := &authCode{
		CodeID:      code.CodeID,
		ExpiresIn:   code.ExpiresIn,
		ClientID:    code.ClientID,
		RedirectURL: code.RedirectURL,
		UserID:      code.UserID,
	}

	col := h.dbClient.Database(databaseName).Collection(authCodeCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert auth code to mongodb")
	}

	return nil
}

// Delete ...
func (h *AuthCodeHandler) Delete(codeID string) error {
	col := h.dbClient.Database(databaseName).Collection(authCodeCollectionName)
	filter := bson.D{
		{Key: "codeID", Value: codeID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete auth code from mongodb")
	}
	return nil
}

// Get ...
func (h *AuthCodeHandler) Get(codeID string) (*model.AuthCode, error) {
	col := h.dbClient.Database(databaseName).Collection(authCodeCollectionName)
	filter := bson.D{
		{Key: "codeID", Value: codeID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &model.AuthCode{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Cause(model.ErrNoSuchCode)
		}
		return nil, errors.Wrap(err, "Failed to get auth code from mongodb")
	}

	return res, nil
}
