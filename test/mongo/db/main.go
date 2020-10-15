package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/db/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// TODO set from arg
	mongoAddr := "mongodb://localhost:27017"
	dbName := "hekate-test"

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

	// Test DB Methods
	testProject(dbClient)
	testClient(dbClient)
	testCustomRole(dbClient)
	testLoginSession(dbClient)
	testSession(dbClient)
	testUser(dbClient)

	// TODO Test Audit Methods

	fmt.Println("Successfully finished")
}

func testProject(dbClient *mongodriver.Client) {
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
		errors.Print(errors.Append(err, "Failed to add master project"))
		os.Exit(1)
	}
	ent.Name = "project1"
	projectHandler.Add(ent)

	ent.UserLock.Enabled = true
	if err := projectHandler.Update(ent); err != nil {
		errors.Print(errors.Append(err, "Failed to update project"))
		os.Exit(1)
	}

	prjs, err := projectHandler.GetList(nil)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get project list"))
		os.Exit(1)
	}
	if len(prjs) != 2 {
		errors.Print(errors.New("", "Unexpected project num was specified. want: 2, got: %d", len(prjs)))
		os.Exit(1)
	}
	for _, prj := range prjs {
		if prj.Name == "project1" {
			if !prj.UserLock.Enabled {
				errors.Print(errors.New("", "Project update failed. project.UserLock.Enabled should be true"))
				os.Exit(1)
			}
		}
	}

	if err := projectHandler.Delete("project1"); err != nil {
		errors.Print(errors.Append(err, "Failed to delete project"))
		os.Exit(1)
	}
}

func testClient(dbClient *mongodriver.Client) {
	clientHandler, err := mongo.NewClientHandler(dbClient)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to create client handler"))
		os.Exit(1)
	}

	ent := &model.ClientInfo{
		ID:          "test-client",
		ProjectName: "master",
		AccessType:  "public",
		CreatedAt:   time.Now(),
	}
	if err := clientHandler.Add("master", ent); err != nil {
		errors.Print(errors.Append(err, "Failed to add client"))
		os.Exit(1)
	}
	ent.ID = "client-2"
	clientHandler.Add("master", ent)
	ent.ID = "client-3"
	clientHandler.Add("master", ent)

	ent.Secret = "mysecret"
	if err := clientHandler.Update("master", ent); err != nil {
		errors.Print(errors.Append(err, "Failed to update client"))
		os.Exit(1)
	}

	clis, err := clientHandler.GetList("master", nil)
	if err != nil {
		errors.Print(errors.Append(err, "Failed to get client list"))
		os.Exit(1)
	}
	if len(clis) != 3 {
		errors.Print(errors.New("", "Unexpected client num was specified. want: 3, got: %d", len(clis)))
		os.Exit(1)
	}
	for _, cli := range clis {
		if cli.ID == "client-3" {
			if ent.Secret != "mysecret" {
				errors.Print(errors.New("", "Client update failed. client.Secret should be mysecret"))
				os.Exit(1)
			}
		}
	}

	if err := clientHandler.Delete("master", "client-3"); err != nil {
		errors.Print(errors.Append(err, "Failed to delete client"))
		os.Exit(1)
	}
	clis, _ = clientHandler.GetList("master", nil)
	if len(clis) != 2 {
		errors.Print(errors.New("", "Unexpected client num was specified after delete. want: 2, got: %d", len(clis)))
		os.Exit(1)
	}

	if err := clientHandler.DeleteAll("master"); err != nil {
		errors.Print(errors.Append(err, "Failed to delete all client"))
		os.Exit(1)
	}
	clis, _ = clientHandler.GetList("master", nil)
	if len(clis) != 0 {
		errors.Print(errors.New("", "Unexpected client num was specified after all delete. want: 0, got: %d", len(clis)))
		os.Exit(1)
	}
}

func testCustomRole(dbClient *mongodriver.Client) {
	// TODO Custom Role Add, Delete, GetList, Update, DeleteAll
}

func testLoginSession(dbClient *mongodriver.Client) {
	// TODO Login Session Add, Delete, GetByCode, Update, DeleteAll, Get, Cleanup
}

func testSession(dbClient *mongodriver.Client) {
	// TODO Session Add, Delete, GetList, DeleteAll
}

func testUser(dbClient *mongodriver.Client) {
	// TODO User Add, Delete, GetList, Update, DeleteAll, AddRole, DeleteRole, DeleteAllCustomRole
}
