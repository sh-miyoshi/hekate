package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoginSessionHandler implement db.LoginSessionHandler
type LoginSessionHandler struct {
	dbClient *mongo.Client
}

// NewLoginSessionHandler ...
func NewLoginSessionHandler(dbClient *mongo.Client) (*LoginSessionHandler, *errors.Error) {
	res := &LoginSessionHandler{
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
	if err != nil {
		return nil, errors.New("", "Failed to create index: %v", err)
	}

	return res, nil
}

// Add ...
func (h *LoginSessionHandler) Add(projectName string, ent *model.LoginSession) *errors.Error {
	v := &loginSession{
		SessionID:    ent.SessionID,
		Code:         ent.Code,
		ExpiresIn:    ent.ExpiresIn,
		Scope:        ent.Scope,
		ResponseType: ent.ResponseType,
		ClientID:     ent.ClientID,
		RedirectURI:  ent.RedirectURI,
		Nonce:        ent.Nonce,
		ProjectName:  ent.ProjectName,
		ResponseMode: ent.ResponseMode,
		Prompt:       ent.Prompt,
		LoginDate:    ent.LoginDate,
	}

	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("", "Failed to insert login session to mongodb: %v", err)
	}

	return nil
}

// Update ...
func (h *LoginSessionHandler) Update(projectName string, ent *model.LoginSession) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authCodeCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: ent.SessionID},
	}

	v := &loginSession{
		SessionID:    ent.SessionID,
		Code:         ent.Code,
		ExpiresIn:    ent.ExpiresIn,
		Scope:        ent.Scope,
		ResponseType: ent.ResponseType,
		ClientID:     ent.ClientID,
		RedirectURI:  ent.RedirectURI,
		Nonce:        ent.Nonce,
		ProjectName:  ent.ProjectName,
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
		return errors.New("", "Failed to update auth codoe session in mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *LoginSessionHandler) Delete(projectName string, sessionID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: sessionID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("", "Failed to delete login session from mongodb: %v", err)
	}
	return nil
}

// GetByCode ...
func (h *LoginSessionHandler) GetByCode(projectName string, code string) (*model.LoginSession, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "code", Value: code},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &loginSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchLoginSession
		}
		return nil, errors.New("", "Failed to get login session from mongodb: %v", err)
	}

	return &model.LoginSession{
		SessionID:    res.SessionID,
		Code:         res.Code,
		ExpiresIn:    res.ExpiresIn,
		Scope:        res.Scope,
		ResponseType: res.ResponseType,
		ClientID:     res.ClientID,
		RedirectURI:  res.RedirectURI,
		Nonce:        res.Nonce,
		ProjectName:  res.ProjectName,
		ResponseMode: res.ResponseMode,
		Prompt:       res.Prompt,
		LoginDate:    res.LoginDate,
	}, nil
}

// Get ...
func (h *LoginSessionHandler) Get(projectName string, id string) (*model.LoginSession, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "sessionID", Value: id},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &loginSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchLoginSession
		}
		return nil, errors.New("", "Failed to get login session from mongodb: %v", err)
	}

	return &model.LoginSession{
		SessionID:    res.SessionID,
		Code:         res.Code,
		ExpiresIn:    res.ExpiresIn,
		Scope:        res.Scope,
		ResponseType: res.ResponseType,
		ClientID:     res.ClientID,
		RedirectURI:  res.RedirectURI,
		Nonce:        res.Nonce,
		ProjectName:  res.ProjectName,
		ResponseMode: res.ResponseMode,
		Prompt:       res.Prompt,
		LoginDate:    res.LoginDate,
	}, nil
}

// DeleteAllInClient ...
func (h *LoginSessionHandler) DeleteAllInClient(projectName string, clientID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "clientID", Value: clientID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("", "Failed to delete authcode session from mongodb: %v", err)
	}
	return nil
}

// DeleteAllInUser ...
func (h *LoginSessionHandler) DeleteAllInUser(projectName string, userID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
		{Key: "userID", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("", "Failed to delete authcode session from mongodb: %v", err)
	}
	return nil
}

// DeleteAllInProject ...
func (h *LoginSessionHandler) DeleteAllInProject(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("", "Failed to delete authcode session from mongodb: %v", err)
	}
	return nil
}
