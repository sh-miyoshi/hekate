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

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	dbClient *mongo.Client
}

// NewUserHandler ...
func NewUserHandler(dbClient *mongo.Client) (*UserInfoHandler, *errors.Error) {
	res := &UserInfoHandler{
		dbClient: dbClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(userCollectionName)
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
		logger.Info("Create index for user")
		// Create Index to Project Name and User ID
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
func (h *UserInfoHandler) Add(projectName string, ent *model.UserInfo) *errors.Error {
	usr := &userInfo{
		ID:           ent.ID,
		ProjectName:  ent.ProjectName,
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PasswordHash: ent.PasswordHash,
		SystemRoles:  ent.SystemRoles,
		CustomRoles:  ent.CustomRoles,
		LockState: lockState{
			Locked:            ent.LockState.Locked,
			VerifyFailedTimes: ent.LockState.VerifyFailedTimes,
		},
	}

	uroles := []interface{}{}
	for _, r := range ent.CustomRoles {
		uroles = append(uroles, &customRoleInUser{
			ProjectName:  ent.ProjectName,
			UserID:       ent.ID,
			CustomRoleID: r,
		})
	}

	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.InsertOne(ctx, usr); err != nil {
		return errors.New("DB failed", "Failed to insert user to mongodb: %v", err)
	}

	if len(uroles) > 0 {
		rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)

		if _, err := rcol.InsertMany(ctx, uroles); err != nil {
			return errors.New("DB failed", "Failed to insert role info in user to mongodb: %v", err)
		}
	}

	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.DeleteOne(ctx, filter); err != nil {
		return errors.New("DB failed", "Failed to delete user from mongodb: %v", err)
	}

	rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)
	filter = bson.D{
		{Key: "user_id", Value: userID},
	}

	if _, err := rcol.DeleteOne(ctx, filter); err != nil {
		return errors.New("DB failed", "Failed to delete custom role in user from mongodb: %v", err)
	}

	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string, filter *model.UserFilter) ([]*model.UserInfo, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)

	f := bson.D{
		{Key: "project_name", Value: projectName},
	}

	if filter != nil {
		if filter.Name != "" {
			f = append(f, bson.E{Key: "name", Value: filter.Name})
		}
		if filter.ID != "" {
			f = append(f, bson.E{Key: "id", Value: filter.ID})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get user list from mongodb: %v", err)
	}

	users := []userInfo{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, errors.New("DB failed", "Failed to parse user list from mongodb: %v", err)
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
			LockState: model.LockState{
				Locked:            user.LockState.Locked,
				VerifyFailedTimes: user.LockState.VerifyFailedTimes,
			},
		})
	}

	return res, nil
}

// Update ...
func (h *UserInfoHandler) Update(projectName string, ent *model.UserInfo) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
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
		LockState: lockState{
			Locked:            ent.LockState.Locked,
			VerifyFailedTimes: ent.LockState.VerifyFailedTimes,
		},
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.New("DB failed", "Failed to update user in mongodb: %v", err)
	}

	rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)
	filter = bson.D{
		{Key: "user_id", Value: ent.ID},
	}
	if _, err := rcol.DeleteMany(ctx, filter); err != nil {
		return errors.New("DB failed", "Failed to delete previous custom role in user from mongodb: %v", err)
	}
	uroles := []interface{}{}
	for _, r := range ent.CustomRoles {
		uroles = append(uroles, &customRoleInUser{
			ProjectName:  ent.ProjectName,
			UserID:       ent.ID,
			CustomRoleID: r,
		})
	}
	if len(uroles) > 0 {
		if _, err := rcol.InsertMany(ctx, uroles); err != nil {
			return errors.New("DB failed", "Failed to insert role info in user to mongodb: %v", err)
		}
	}

	return nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete user from mongodb: %v", err)
	}

	rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)
	if _, err := rcol.DeleteMany(ctx, filter); err != nil {
		return errors.New("DB failed", "Failed to delete custom role in user from mongodb: %v", err)
	}

	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(projectName string, userID string, roleType model.RoleType, roleID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	user := &userInfo{}
	if err := col.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.ErrNoSuchUser
		}
		return errors.New("DB failed", "Failed to get user from mongodb: %v", err)
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
		return errors.New("DB failed", "Failed to add role to user in mongodb: %v", err)
	}

	if roleType == model.RoleCustom {
		role := customRoleInUser{
			ProjectName:  user.ProjectName,
			UserID:       user.ID,
			CustomRoleID: roleID,
		}
		rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)
		if _, err := rcol.InsertOne(ctx, role); err != nil {
			return errors.New("DB failed", "Failed to insert role info in user to mongodb: %v", err)
		}
	}

	return nil
}

// DeleteRole ...
func (h *UserInfoHandler) DeleteRole(projectName string, userID string, roleID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "id", Value: userID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	user := &userInfo{}
	if err := col.FindOne(ctx, filter).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return model.ErrNoSuchUser
		}
		return errors.New("DB failed", "Failed to get user from mongodb: %v", err)
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
		return errors.New("DB failed", "Failed to add role to user in mongodb: %v", err)
	}

	rcol := h.dbClient.Database(databaseName).Collection(roleInUserCollectionName)
	filter = bson.D{
		{Key: "user_id", Value: userID},
		{Key: "custom_role_id", Value: roleID},
	}

	if _, err := rcol.DeleteOne(ctx, filter); err != nil {
		return errors.New("DB failed", "Failed to delete custom role in user from mongodb: %v", err)
	}

	return nil
}

// DeleteAllCustomRole ...
func (h *UserInfoHandler) DeleteAllCustomRole(projectName string, roleID string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(userCollectionName)
	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "custom_role_id", Value: roleID},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to get role list from mongodb: %v", err)
	}

	roles := []customRoleInUser{}
	if err := cursor.All(ctx, &roles); err != nil {
		return errors.New("DB failed", "Failed to parse role list from mongodb: %v", err)
	}

	for _, r := range roles {
		h.DeleteRole(projectName, r.UserID, roleID)
	}

	return nil
}
