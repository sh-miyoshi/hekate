package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
)

func main() {
	// TODO set from arg
	mongoAddr := "mongodb://localhost:27017"
	dbName := "hekate-test"

	// TODO log output

	mongo.ChangeDatabase(dbName)
	dbClient, err := mongo.NewClient(mongoAddr)
	if err != nil {
		errors.Print(err)
		os.Exit(1)
	}

	// init database data
	if err := dbClient.Database(dbName).Drop(context.TODO()); err != nil {
		fmt.Printf("Failed to initialize test database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbClient.Database(dbName).Drop(context.TODO()); err != nil {
			fmt.Printf("Failed to cleanup test database: %v", err)
		}
	}()

	// Test Project Methods
	projectHandler, err := mongo.NewProjectHandler(dbClient)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to create project handler"))
		os.Exit(1)
	}
	ent := &model.ProjectInfo{
		Name:        "master",
		CreatedAt:   time.Now(),
		TokenConfig: &model.TokenConfig{},
	}
	if err := projectHandler.Add(ent); err != nil {
		errors.Print(errors.Append(err, "Failed to add new project"))
		os.Exit(1)
	}
	// TODO Project Delete, GetList, Update

	// TODO Client Add, Delete, GetList, Update, DeleteAll
	// TODO Custom Role Add, Delete, GetList, Update, DeleteAll
	// TODO Login Session Add, Delete, GetByCode, Update, DeleteAll, Get, Cleanup
	// TODO Session Add, Delete, GetList, DeleteAll
	// TODO User Add, Delete, GetList, Update, DeleteAll, AddRole, DeleteRole, DeleteAllCustomRole

	fmt.Println("Successfully finished")
}
