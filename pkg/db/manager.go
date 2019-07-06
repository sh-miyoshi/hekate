package db

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

// DBHandler is an interface for handler of db
type DBHandler interface {
	ConnectDB(connectString string) error

	CreateUser(newUser UserRequest) error
	DeleteUser(userID string) error
	GetUserList() ([]User, error)
	UpdatePassowrd(newPassword string) error

	AddRoleToUser(role RoleType, userID string) error
	RemoveRoleFromUser(role RoleType, userID string) error

	SetTokenConfig(config TokenConfig) error
	GetTokenConfig() (TokenConfig, error)
}

var instance DBHandler

// InitDBHandler initialize handler for user
func InitDBHandler(dbType DBType) error {
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
func GetInst() DBHandler {
	return instance
}
