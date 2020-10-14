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

// LoginSessionHandler implement db.LoginSessionHandler
type LoginSessionHandler struct {
	dbClient *mongo.Client
}

// NewLoginSessionHandler ...
func NewLoginSessionHandler(dbClient *mongo.Client) (*LoginSessionHandler, *errors.Error) {
	res := &LoginSessionHandler{
		dbClient: dbClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
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
		logger.Info("Create index for login session")
		// Create Index to Project Name and Login Session ID
		mod := mongo.IndexModel{
			Keys: bson.M{
				"project_name": 1, // index in ascending order
				"id":           1, // index in ascending order
			},
		}
		if _, err := iv.CreateOne(ctx, mod); err != nil {
			return nil, errors.New("DB failed", "Failed to create index: %v", err)
		}
	}

	return res, nil
}

// Add ...
func (h *LoginSessionHandler) Add(projectName string, ent *model.LoginSession) *errors.Error {
	v := &loginSession{
		SessionID:     ent.SessionID,
		Code:          ent.Code,
		ExpiresIn:     ent.ExpiresIn,
		UnixExpiresIn: ent.ExpiresIn.Unix(),
		Scope:         ent.Scope,
		ResponseType:  ent.ResponseType,
		ClientID:      ent.ClientID,
		RedirectURI:   ent.RedirectURI,
		Nonce:         ent.Nonce,
		ProjectName:   ent.ProjectName,
		ResponseMode:  ent.ResponseMode,
		Prompt:        ent.Prompt,
		UserID:        ent.UserID,
		LoginDate:     ent.LoginDate,
	}

	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("DB failed", "Failed to insert login session to mongodb: %v", err)
	}

	return nil
}

// Update ...
func (h *LoginSessionHandler) Update(projectName string, ent *model.LoginSession) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "session_id", Value: ent.SessionID},
	}

	v := &loginSession{
		SessionID:     ent.SessionID,
		Code:          ent.Code,
		ExpiresIn:     ent.ExpiresIn,
		UnixExpiresIn: ent.ExpiresIn.Unix(),
		Scope:         ent.Scope,
		ResponseType:  ent.ResponseType,
		ClientID:      ent.ClientID,
		RedirectURI:   ent.RedirectURI,
		Nonce:         ent.Nonce,
		ProjectName:   ent.ProjectName,
		ResponseMode:  ent.ResponseMode,
		Prompt:        ent.Prompt,
		UserID:        ent.UserID,
		LoginDate:     ent.LoginDate,
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.New("DB failed", "Failed to update auth codoe session in mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *LoginSessionHandler) Delete(projectName string, filter *model.LoginSessionFilter) *errors.Error {
	if filter == nil {
		return nil
	}

	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	f := bson.D{
		{Key: "project_name", Value: projectName},
	}

	if filter != nil {
		if filter.SessionID != "" {
			f = append(f, bson.E{Key: "session_id", Value: filter.SessionID})
		}
		if filter.UserID != "" {
			f = append(f, bson.E{Key: "user_id", Value: filter.UserID})
		}
		if filter.ClientID != "" {
			f = append(f, bson.E{Key: "client_id", Value: filter.ClientID})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, f)
	if err != nil {
		return errors.New("DB failed", "Failed to delete login session from mongodb: %v", err)
	}
	return nil
}

// GetByCode ...
func (h *LoginSessionHandler) GetByCode(projectName string, code string) (*model.LoginSession, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "code", Value: code},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &loginSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchLoginSession
		}
		return nil, errors.New("DB failed", "Failed to get login session from mongodb: %v", err)
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
		UserID:       res.UserID,
		LoginDate:    res.LoginDate,
	}, nil
}

// Get ...
func (h *LoginSessionHandler) Get(projectName string, id string) (*model.LoginSession, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "session_id", Value: id},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &loginSession{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchLoginSession
		}
		return nil, errors.New("DB failed", "Failed to get login session from mongodb: %v", err)
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
		UserID:       res.UserID,
		LoginDate:    res.LoginDate,
	}, nil
}

// DeleteAll ...
func (h *LoginSessionHandler) DeleteAll(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
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
func (h *LoginSessionHandler) Cleanup(now time.Time) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(authcodeSessionCollectionName)
	filter := bson.D{
		{Key: "unix_expires_in", Value: bson.D{{Key: "$lt", Value: now.Unix()}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete expired session from mongodb: %v", err)
	}

	return nil
}
