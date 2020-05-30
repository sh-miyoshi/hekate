package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
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
func NewProjectHandler(dbClient *mongo.Client) (*ProjectInfoHandler, error) {
	res := &ProjectInfoHandler{
		dbClient: dbClient,
	}

	// Create Index to Project Name
	mod := mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	col := res.dbClient.Database(databaseName).Collection(projectCollectionName)
	_, err := col.Indexes().CreateOne(ctx, mod)

	return res, err
}

// Add ...
func (h *ProjectInfoHandler) Add(ent *model.ProjectInfo) error {
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
	}
	for _, t := range ent.AllowGrantTypes {
		v.AllowGrantTypes = append(v.AllowGrantTypes, t.String())
	}

	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.Wrap(err, "Failed to insert project to mongodb")
	}

	return nil
}

// Delete ...
func (h *ProjectInfoHandler) Delete(name string) error {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return errors.Wrap(err, "Failed to delete project from mongodb")
	}
	return nil
}

// GetList ...
func (h *ProjectInfoHandler) GetList() ([]*model.ProjectInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get project list from mongodb")
	}

	projects := []projectInfo{}
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, errors.Wrap(err, "Failed to get project list from mongodb")
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
		}
		for _, t := range prj.AllowGrantTypes {
			typ, _ := model.GetGrantType(t)
			info.AllowGrantTypes = append(info.AllowGrantTypes, typ)
		}
		res = append(res, info)
	}

	return res, nil
}

// Get ...
func (h *ProjectInfoHandler) Get(name string) (*model.ProjectInfo, error) {
	col := h.dbClient.Database(databaseName).Collection(projectCollectionName)
	filter := bson.D{
		{Key: "name", Value: name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	project := &projectInfo{}
	if err := col.FindOne(ctx, filter).Decode(project); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, model.ErrNoSuchProject
		}
		return nil, errors.Wrap(err, "Failed to get project from mongodb")
	}
	logger.Debug("Get project %s data: %v", name, project)

	res := &model.ProjectInfo{
		Name:         project.Name,
		CreatedAt:    project.CreatedAt,
		PermitDelete: project.PermitDelete,
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  project.TokenConfig.AccessTokenLifeSpan,
			RefreshTokenLifeSpan: project.TokenConfig.RefreshTokenLifeSpan,
			SigningAlgorithm:     project.TokenConfig.SigningAlgorithm,
			SignPublicKey:        project.TokenConfig.SignPublicKey,
			SignSecretKey:        project.TokenConfig.SignSecretKey,
		},
		PasswordPolicy: model.PasswordPolicy{
			MinimumLength:       project.PasswordPolicy.MinimumLength,
			NotUserName:         project.PasswordPolicy.NotUserName,
			BlackList:           project.PasswordPolicy.BlackList,
			UseCharacter:        model.CharacterType(project.PasswordPolicy.UseCharacter),
			UseDigit:            project.PasswordPolicy.UseDigit,
			UseSpecialCharacter: project.PasswordPolicy.UseSpecialCharacter,
		},
	}
	for _, t := range project.AllowGrantTypes {
		typ, _ := model.GetGrantType(t)
		res.AllowGrantTypes = append(res.AllowGrantTypes, typ)
	}

	return res, nil
}

// Update ...
func (h *ProjectInfoHandler) Update(ent *model.ProjectInfo) error {
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
	}
	for _, t := range ent.AllowGrantTypes {
		v.AllowGrantTypes = append(v.AllowGrantTypes, t.String())
	}

	updates := bson.D{
		{Key: "$set", Value: v},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if _, err := col.UpdateOne(ctx, filter, updates); err != nil {
		return errors.Wrap(err, "Failed to update project in mongodb")
	}

	return nil
}
