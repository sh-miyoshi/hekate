package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sh-miyoshi/hekate/pkg/audit/model"
	"github.com/sh-miyoshi/hekate/pkg/audit/mongo"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

func main() {
	dbName := "hekate-audit-test"
	mongoAddr := flag.String("mongo-addr", "mongodb://localhost:27017", "connection string for mongodb")
	flag.Parse()

	mongo.ChangeDatabase(dbName)
	dbClient, err := mongo.NewClient(*mongoAddr)
	if err != nil {
		fmt.Printf("Failed to create mongo client: %v\n", err)
		os.Exit(1)
	}

	// init database data
	if err := dbClient.Database(dbName).Drop(context.TODO()); err != nil {
		fmt.Printf("Failed to initialize test database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbClient.Database(dbName).Drop(context.TODO()); err != nil {
			fmt.Printf("Failed to cleanup test database: %v\n", err)
		}
	}()

	// Test Audit Methods
	handler := mongo.NewHandler(dbClient)
	if err := handler.Save("master", time.Now(), "USER", "POST", "/user", ""); err != nil {
		errors.Print(err)
		os.Exit(1)
	}
	// Set test data
	baseTime := time.Now()
	for i := -50; i < 250; i++ {
		tm := baseTime.AddDate(0, 0, i)
		if err := handler.Save("project1", tm, "TEST", "GET", "/", ""); err != nil {
			errors.Print(err)
			os.Exit(1)
		}
	}

	// Test fromDate and toDate
	fromDate := util.TimeTruncate(baseTime)
	n := 10
	toDate := util.TimeTruncate(baseTime.AddDate(0, 0, n-1))
	res, _ := handler.Get("project1", fromDate, toDate, 0)
	if len(res) != n {
		fmt.Printf("Failed to get audit events by filtering date. expect: %d, but got: %d\n", n, len(res))
		os.Exit(1)
	}

	// Test max num
	fromDate = util.TimeTruncate(baseTime.AddDate(0, 0, -50))
	toDate = util.TimeTruncate(baseTime.AddDate(0, 0, 250))
	res, _ = handler.Get("project1", fromDate, toDate, 0)
	if len(res) != model.AuditGetMaxNum {
		fmt.Printf("Failed to get audit events by max num. expect: %d, but got: %d\n", model.AuditGetMaxNum, len(res))
		os.Exit(1)
	}

	// Test offset
	fromDate = util.TimeTruncate(baseTime)
	offset := uint(1)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	tm100 := baseTime.AddDate(0, 0, 100)
	if !tmEqual(res[0].Time, tm100) {
		fmt.Printf("Failed to get audit event by offset. expect: %v, but got: %v\n", tm100, res[0].Time)
		os.Exit(1)
	}
	offset = uint(2)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	tm200 := baseTime.AddDate(0, 0, 200)
	if !tmEqual(res[0].Time, tm200) {
		fmt.Printf("Failed to get audit event by offset 2. expect: %v, but got: %v\n", tm200, res[0].Time)
		os.Exit(1)
	}
	if len(res) != 50 {
		fmt.Printf("Failed to get audit event num by offset 2. expect: %d, but got: %d\n", 50, len(res))
		os.Exit(1)
	}
	offset = uint(100)
	res, _ = handler.Get("project1", fromDate, toDate, offset)
	if len(res) != 0 {
		fmt.Printf("Failed to get audit event by over offset. expect: 0, but got: %d\n", len(res))
		os.Exit(1)
	}

	fmt.Println("Successfully finished")
}

func tmEqual(tm1, tm2 time.Time) bool {
	tm1 = util.TimeTruncate(tm1.UTC())
	tm2 = util.TimeTruncate(tm2.UTC())
	return tm1.Equal(tm2)
}
