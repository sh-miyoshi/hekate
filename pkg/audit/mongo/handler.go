package mongo

import (
	"context"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "audit"
	timeoutSecond  = 5
)

var (
	databaseName = "hekate"
)

// Handler ...
type Handler struct {
	dbClient *mongo.Client
}

// NewClient ...
func NewClient(connStr string) (*mongo.Client, *errors.Error) {
	cli, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, errors.New("DB failed", "Failed to create mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := cli.Connect(ctx); err != nil {
		return nil, errors.New("DB failed", "Failed to connect to mongo: %v", err)
	}

	if err := cli.Ping(ctx, nil); err != nil {
		return nil, errors.New("DB failed", "Failed to ping to mongo: %v", err)
	}

	return cli, nil
}

// NewHandler ...
func NewHandler(dbClient *mongo.Client) *Handler {
	return &Handler{
		dbClient: dbClient,
	}
}

// ChangeDatabase changes the database of store data
// this method should call at first
func ChangeDatabase(name string) {
	if len(name) > 0 {
		databaseName = name
	}
}

// Ping ...
func (h *Handler) Ping() *errors.Error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	if err := h.dbClient.Ping(ctx, nil); err != nil {
		return errors.New("DB failed", "DB Ping failed: %v", err)
	}
	return nil
}

// Save ...
func (h *Handler) Save(projectName string, tm time.Time, resType, method, path, message string) *errors.Error {
	v := &audit{
		ProjectName:  projectName,
		Time:         tm,
		ResourceType: resType,
		Method:       method,
		Path:         path,
		IsSuccess:    message == "",
		Message:      message,
		UnixTime:     tm.Unix(),
	}

	col := h.dbClient.Database(databaseName).Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	_, err := col.InsertOne(ctx, v)
	if err != nil {
		return errors.New("DB failed", "Failed to insert audit to mongodb: %v", err)
	}

	return nil
}

// Get ...
func (h *Handler) Get(projectName string, fromDate, toDate time.Time, offset uint) ([]model.Audit, *errors.Error) {
	// TODO(implement this)
	col := h.dbClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{
		{Key: "project_name", Value: projectName},
		{Key: "unixtime", Value: bson.D{{Key: "$gte", Value: fromDate.Unix()}}},
	}

	// if we want to get logs whose date are from "2019-09-19",
	// we have to pass "2019-09-20 00:00:00.000" to mongodb.
	toDate = toDate.AddDate(0, 0, 1)
	toDate = toDate.Add(-time.Nanosecond)
	filter = append(filter, bson.E{Key: "unixtime", Value: bson.D{{Key: "$lt", Value: toDate.Unix()}}})

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "unixtime", Value: -1}})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(model.AuditGetMaxNum))

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()

	cursor, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, errors.New("DB failed", "Failed to get audit list from mongodb: %v", err)
	}

	audits := []audit{}
	if err := cursor.All(ctx, &audits); err != nil {
		return nil, errors.New("DB failed", "Failed to get audit list from mongodb: %v", err)
	}

	res := []model.Audit{}
	for _, a := range audits {
		res = append(res, model.Audit{
			ProjectName:  a.ProjectName,
			Time:         a.Time,
			ResourceType: a.ResourceType,
			Method:       a.Method,
			Path:         a.Path,
			IsSuccess:    a.IsSuccess,
			Message:      a.Message,
		})
	}

	return res, nil
}
