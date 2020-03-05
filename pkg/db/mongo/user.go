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

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	session  mongo.Session
	dbClient *mongo.Client
}

// NewUserHandler ...
func NewUserHandler(dbClient *mongo.Client) (*UserInfoHandler, error) {
	res := &UserInfoHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name
	mod := mongo.IndexModel{
		Keys: bson.M{
			"id": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(userCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	v := &userInfo{
		ID:           ent.ID,
		ProjectName:  ent.ProjectName,
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PasswordHash: ent.PasswordHash,
		SystemRoles:  ent.SystemRoles,
		CustomRoles:  ent.CustomRoles,
	}

	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert user to mongodb")
	}

	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(userID string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete user from mongodb")
	}
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	f := bson.D{
		{Key: "projectName", Value: projectName},
	}

	if filter != nil {
		if filter.Name != "" {
			f = append(f, bson.E{Key: "name", Value: filter.Name})
		}
		// TODO(add other filter)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user list from mongodb")
	}

	users := []userInfo{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, errors.Wrap(err, "Failed to get user list from mongodb")
	}

	res := []*model.UserInfo{}
	for _, user := range users {
		res = append(res, &model.UserInfo{
			ID:           user.ID,
			ProjectName:  user.ProjectName,
			Name:         user.Name,
			CreatedAt:    user.CreatedAt,
			PasswordHash: user.PasswordHash,
			SystemRoles:  user.SystemRoles,
			CustomRoles:  user.CustomRoles,
		})
	}

	return res, nil
}

// Get ...
func (h *UserInfoHandler) Get(userID string) (*model.UserInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	res := &userInfo{}
	if err := col.FindOne(ctx, filter).Decode(res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchUser
		}
		return nil, errors.Wrap(err, "Failed to get user from mongodb")
	}

	return &model.UserInfo{
		ID:           res.ID,
		ProjectName:  res.ProjectName,
		Name:         res.Name,
		CreatedAt:    res.CreatedAt,
		PasswordHash: res.PasswordHash,
		SystemRoles:  res.SystemRoles,
		CustomRoles:  res.CustomRoles,
	}, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "id", Value: ent.ID},
	}

	v := &userInfo{
		ID:           ent.ID,
		ProjectName:  ent.ProjectName,
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PasswordHash: ent.PasswordHash,
		SystemRoles:  ent.SystemRoles,
		CustomRoles:  ent.CustomRoles,
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to update user in mongodb")
	}

	return nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "projectName", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete user from mongodb")
	}
	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(userID string, roleType model.RoleType, roleID string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	user := &userInfo{}
	if err := col.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.ErrNoSuchUser
		}
		return errors.Wrap(err, "Failed to get user from mongodb")
	}

	if roleType == model.RoleSystem {
		for _, role := range user.SystemRoles {
			if role == roleID {
				return model.ErrRoleAlreadyAppended
			}
		}
		user.SystemRoles = append(user.SystemRoles, roleID)
	} else if roleType == model.RoleCustom {
		for _, role := range user.CustomRoles {
			if role == roleID {
				return model.ErrRoleAlreadyAppended
			}
		}
		user.CustomRoles = append(user.CustomRoles, roleID)
	}

	updates := bson.D{
		{Key: "$set", Value: user},
	}

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to add role to user in mongodb")
	}

	return nil
}

// DeleteRole ...
func (h *UserInfoHandler) DeleteRole(userID string, roleID string) error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	user := &userInfo{}
	if err := col.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.ErrNoSuchUser
		}
		return errors.Wrap(err, "Failed to get user from mongodb")
	}

	deleted := false
	roles := []string{}
	for _, role := range user.SystemRoles {
		if role == roleID {
			deleted = true
		} else {
			roles = append(roles, role)
		}
	}

	if deleted {
		user.SystemRoles = roles
	} else {
		deleted = false
		roles = []string{}
		for _, role := range user.CustomRoles {
			if role == roleID {
				deleted = true
			} else {
				roles = append(roles, role)
			}
		}
		if !deleted {
			return model.ErrNoSuchRoleInUser
		}
		user.CustomRoles = roles
	}

	updates := bson.D{
		{Key: "$set", Value: user},
	}

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to add role to user in mongodb")
	}

	return nil
}

// BeginTx ...
func (h *UserInfoHandler) BeginTx() error {
	var err error
	h.session, err = h.dbClient.StartSession()
	if err != nil {
		return err
	}
	err = h.session.StartTransaction()
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
		defer cancel()
		h.session.EndSession(ctx)
		return err
	}
	return nil
}

// CommitTx ...
func (h *UserInfoHandler) CommitTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.CommitTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}

// AbortTx ...
func (h *UserInfoHandler) AbortTx() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	err := h.session.AbortTransaction(ctx)
	h.session.EndSession(ctx)
	return err
}
