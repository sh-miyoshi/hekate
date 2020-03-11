package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoginSessionHandler implement db.LoginSessionHandler
type LoginSessionHandler struct {
	dbClient *mongo.Client
}

// NewLoginSessionHandler ...
func NewLoginSessionHandler(dbClient *mongo.Client) (*LoginSessionHandler, error) {
	res := &LoginSessionHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name
	mod := mongo.IndexModel{
		Keys: bson.M{
			"verifyCode": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(loginSessionCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *LoginSessionHandler) Add(info *model.LoginSessionInfo) error {
	v := &loginSessionInfo{
		VerifyCode:  info.VerifyCode,
		ExpiresIn:   info.ExpiresIn,
		ClientID:    info.ClientID,
		RedirectURI: info.RedirectURI,
	}

	col := h.dbClient.Database(databaseName).Collection(loginSessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert login session to mongodb")
	}

	return nil
}

// Delete ...
func (h *LoginSessionHandler) Delete(code string) error {
	col := h.dbClient.Database(databaseName).Collection(loginSessionCollectionName)
	filter := bson.D{
		{Key: "verifyCode", Value: code},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete login session from mongodb")
	}
	return nil
}

// Get ...
func (h *LoginSessionHandler) Get(code string) (*model.LoginSessionInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(loginSessionCollectionName)
	filter := bson.D{
		{Key: "verifyCode", Value: code},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &loginSessionInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchLoginSession
		}
		return nil, errors.Wrap(err, "Failed to get login session from mongodb")
	}

	return &model.LoginSessionInfo{
		VerifyCode:  res.VerifyCode,
		ExpiresIn:   res.ExpiresIn,
		ClientID:    res.ClientID,
		RedirectURI: res.RedirectURI,
	}, nil
}
