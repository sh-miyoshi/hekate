package userdb

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

// DBType is type of database
type DBType int

const (
	// DBRemote use remote database
	DBRemote DBType = iota
	// DBLocal use local csv file for database
	DBLocal
)

// UserHandler is an interface for handler of user
type UserHandler interface {
	ConnectDB(connectString string) error
	Authenticate(req UserRequest) (string, error)
	CreateUser(newUser UserRequest) error
	DeleteUser(userName string) error
}

var instance UserHandler

// InitUserHandler initialize handler for user
func InitUserHandler(dbType DBType) error {
	switch dbType {
	case DBRemote:
		logger.Info("Run User DB as Remote Mode")
		return fmt.Errorf("Sorry, not implemented yet")
	case DBLocal:
		logger.Info("Run User DB as Local Mode")
		instance = &localDBHandler{}
		return nil
	}
	return fmt.Errorf("No such database type")
}

// GetInst return a instance of handler
func GetInst() UserHandler {
	return instance
}
