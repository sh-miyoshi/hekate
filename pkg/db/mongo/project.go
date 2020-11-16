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

// ProjectInfoHandler implement db.ProjectInfoHandler
type ProjectInfoHandler struct {
	dbClient *mongo.Client
}

// NewProjectHandler ...
func NewProjectHandler(dbClient *mongo.Client) (*ProjectInfoHandler, *errors.Error) {
	res := &ProjectInfoHandler{
		dbClient: dbClient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	// Get index info
	col := res.dbClient.Database(databaseName).Collection(projectCollectionName)
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
		logger.Info("Create index for project")
		// Create Index to Project Name
		mod := mongo.IndexModel{
			Keys: bson.M{
				"name": 1, // index in ascending order
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
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) *errors.Error {
	v := &projectInfo{
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PermitDelete: ent.PermitDelete,
		TokenConfig: &tokenConfig{
			AccessTokenLifeSpan:  ent.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: ent.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     ent.TokenConfig.SigningAlgorithm,
			SignPublicKey:        ent.TokenConfig.SignPublicKey,
			SignSecretKey:        ent.TokenConfig.SignSecretKey,
		},
		PasswordPolicy: passwordPolicy{
			MinimumLength:       ent.PasswordPolicy.MinimumLength,
			NotUserName:         ent.PasswordPolicy.NotUserName,
			BlackList:           ent.PasswordPolicy.BlackList,
			UseCharacter:        string(ent.PasswordPolicy.UseCharacter),
			UseDigit:            ent.PasswordPolicy.UseDigit,
			UseSpecialCharacter: ent.PasswordPolicy.UseSpecialCharacter,
		},
		UserLock: userLock{
			Enabled:          ent.UserLock.Enabled,
			MaxLoginFailure:  ent.UserLock.MaxLoginFailure,
			LockDuration:     ent.UserLock.LockDuration,
			FailureResetTime: ent.UserLock.FailureResetTime,
		},
	}
	for _, t := range ent.AllowGrantTypes {
		v.AllowGrantTypes = append(v.AllowGrantTypes, string(t))
	}

	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("DB failed", "Failed to insert project to mongodb: %v", err)
	}

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("DB failed", "Failed to delete project from mongodb: %v", err)
	}
	return nil
}

// GetList ...
func (h *ProjectInfoHandler) GetList(filter *model.ProjectFilter) ([]*model.ProjectInfo, *errors.Error) {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	f := bson.D{}
	if filter != nil {
		if filter.Name != "" {
			f = append(f, bson.E{Key: "name", Value: filter.Name})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, f)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get project list from mongodb: %v", err)
	}

	projects := []projectInfo{}
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, errors.New("DB failed", "Failed to parse project list: %v", err)
	}

	res := []*model.ProjectInfo{}
	for _, prj := range projects {
		info := &model.ProjectInfo{
			Name:         prj.Name,
			CreatedAt:    prj.CreatedAt,
			PermitDelete: prj.PermitDelete,
			TokenConfig: &model.TokenConfig{
				AccessTokenLifeSpan:  prj.TokenConfig.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: prj.TokenConfig.RefreshTokenLifeSpan,
				SigningAlgorithm:     prj.TokenConfig.SigningAlgorithm,
				SignPublicKey:        prj.TokenConfig.SignPublicKey,
				SignSecretKey:        prj.TokenConfig.SignSecretKey,
			},
			PasswordPolicy: model.PasswordPolicy{
				MinimumLength:       prj.PasswordPolicy.MinimumLength,
				NotUserName:         prj.PasswordPolicy.NotUserName,
				BlackList:           prj.PasswordPolicy.BlackList,
				UseCharacter:        model.CharacterType(prj.PasswordPolicy.UseCharacter),
				UseDigit:            prj.PasswordPolicy.UseDigit,
				UseSpecialCharacter: prj.PasswordPolicy.UseSpecialCharacter,
			},
			UserLock: model.UserLock{
				Enabled:          prj.UserLock.Enabled,
				MaxLoginFailure:  prj.UserLock.MaxLoginFailure,
				LockDuration:     prj.UserLock.LockDuration,
				FailureResetTime: prj.UserLock.FailureResetTime,
			},
		}
		for _, t := range prj.AllowGrantTypes {
			info.AllowGrantTypes = append(info.AllowGrantTypes, model.GrantType(t))
		}
		res = append(res, info)
	}

	return res, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) *errors.Error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: ent.Name},
	}

	v := &projectInfo{
		Name:         ent.Name,
		CreatedAt:    ent.CreatedAt,
		PermitDelete: ent.PermitDelete,
		TokenConfig: &tokenConfig{
			AccessTokenLifeSpan:  ent.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: ent.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     ent.TokenConfig.SigningAlgorithm,
			SignPublicKey:        ent.TokenConfig.SignPublicKey,
			SignSecretKey:        ent.TokenConfig.SignSecretKey,
		},
		PasswordPolicy: passwordPolicy{
			MinimumLength:       ent.PasswordPolicy.MinimumLength,
			NotUserName:         ent.PasswordPolicy.NotUserName,
			BlackList:           ent.PasswordPolicy.BlackList,
			UseCharacter:        string(ent.PasswordPolicy.UseCharacter),
			UseDigit:            ent.PasswordPolicy.UseDigit,
			UseSpecialCharacter: ent.PasswordPolicy.UseSpecialCharacter,
		},
		UserLock: userLock{
			Enabled:          ent.UserLock.Enabled,
			MaxLoginFailure:  ent.UserLock.MaxLoginFailure,
			LockDuration:     ent.UserLock.LockDuration,
			FailureResetTime: ent.UserLock.FailureResetTime,
		},
	}
	for _, t := range ent.AllowGrantTypes {
		v.AllowGrantTypes = append(v.AllowGrantTypes, string(t))
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.New("DB failed", "Failed to update project in mongodb: %v", err)
	}

	return nil
}
