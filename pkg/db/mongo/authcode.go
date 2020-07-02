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

// AuthCodeSessionHandler implement db.AuthCodeSessionHandler
type AuthCodeSessionHandler struct {
	dbClient *mongo.Client
}

// NewAuthCodeSessionHandler ...
func NewAuthCodeSessionHandler(dbClient *mongo.Client) (*AuthCodeSessionHandler, error) {
	res := &AuthCodeSessionHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name and Session ID
	mod := mongo.IndexModel{
		Keys: bson.M{
			"projectName": 1, // index in ascending order
			"sessionID":   1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *AuthCodeSessionHandler) Add(projectName string, ent *model.AuthCodeSession) error {
	v := &authCodeSession{
		SessionID:    ent.SessionID,
		Code:         ent.Code,
		ExpiresIn:    ent.ExpiresIn,
		Scope:        ent.Scope,
		ResponseType: ent.ResponseType,
		ClientID:     ent.ClientID,
		RedirectURI:  ent.RedirectURI,
		Nonce:        ent.Nonce,
		ProjectName:  ent.ProjectName,
		MaxAge:       ent.MaxAge,
		ResponseMode: ent.ResponseMode,
		Prompt:       ent.Prompt,
		LoginDate:    ent.LoginDate,
	}

	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert login session to mongodb")
	}

	return nil
}

// Update ...
func (h *AuthCodeSessionHandler) Update(projectName string, ent *model.AuthCodeSession) error {
	col := h.dbClient.Database(databaseName).Collection(authCodeCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: ent.SessionID},
	}

	v := &authCodeSession{
		SessionID:    ent.SessionID,
		Code:         ent.Code,
		ExpiresIn:    ent.ExpiresIn,
		Scope:        ent.Scope,
		ResponseType: ent.ResponseType,
		ClientID:     ent.ClientID,
		RedirectURI:  ent.RedirectURI,
		Nonce:        ent.Nonce,
		ProjectName:  ent.ProjectName,
		MaxAge:       ent.MaxAge,
		ResponseMode: ent.ResponseMode,
		Prompt:       ent.Prompt,
		LoginDate:    ent.LoginDate,
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to update auth codoe session in mongodb")
	}

	return nil
}

// Delete ...
func (h *AuthCodeSessionHandler) Delete(projectName string, sessionID string) error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete login session from mongodb")
	}
	return nil
}

// GetByCode ...
func (h *AuthCodeSessionHandler) GetByCode(projectName string, code string) (*model.AuthCodeSession, error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "code", Value: code},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &authCodeSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchAuthCodeSession
		}
		return nil, errors.Wrap(err, "Failed to get login session from mongodb")
	}

	return &model.AuthCodeSession{
		SessionID:    res.SessionID,
		Code:         res.Code,
		ExpiresIn:    res.ExpiresIn,
		Scope:        res.Scope,
		ResponseType: res.ResponseType,
		ClientID:     res.ClientID,
		RedirectURI:  res.RedirectURI,
		Nonce:        res.Nonce,
		ProjectName:  res.ProjectName,
		MaxAge:       res.MaxAge,
		ResponseMode: res.ResponseMode,
		Prompt:       res.Prompt,
		LoginDate:    res.LoginDate,
	}, nil
}

// Get ...
func (h *AuthCodeSessionHandler) Get(projectName string, id string) (*model.AuthCodeSession, error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: id},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &authCodeSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchAuthCodeSession
		}
		return nil, errors.Wrap(err, "Failed to get login session from mongodb")
	}

	return &model.AuthCodeSession{
		SessionID:    res.SessionID,
		Code:         res.Code,
		ExpiresIn:    res.ExpiresIn,
		Scope:        res.Scope,
		ResponseType: res.ResponseType,
		ClientID:     res.ClientID,
		RedirectURI:  res.RedirectURI,
		Nonce:        res.Nonce,
		ProjectName:  res.ProjectName,
		MaxAge:       res.MaxAge,
		ResponseMode: res.ResponseMode,
		Prompt:       res.Prompt,
		LoginDate:    res.LoginDate,
	}, nil
}

// DeleteAllInClient ...
func (h *AuthCodeSessionHandler) DeleteAllInClient(projectName string, clientID string) error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "clientID", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete authcode session from mongodb")
	}
	return nil
}

// DeleteAllInUser ...
func (h *AuthCodeSessionHandler) DeleteAllInUser(projectName string, userID string) error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "userID", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete authcode session from mongodb")
	}
	return nil
}

// DeleteAllInProject ...
func (h *AuthCodeSessionHandler) DeleteAllInProject(projectName string) error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete authcode session from mongodb")
	}
	return nil
}
