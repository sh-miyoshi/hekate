package userdb

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	"io"
	"os"
	"sync"
)

type localDBHandler struct {
	UserHandler

	userFileName string
	mu           sync.Mutex
}

// This func read all csv data at once, so should not use in production
func csvReadAll(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("Failed to open DB file %s in Authenticate: %v", fileName, err)
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'

	return reader.ReadAll()
}

// ConnectDB check file exists(connectString is a file path of user data)
func (l *localDBHandler) ConnectDB(connectString string) error {

	// create new db file
	if _, err := os.Create(connectString); err != nil {
		return err
	}

	// TODO check file broken

	l.userFileName = connectString
	return nil
}

func (l *localDBHandler) Authenticate(req UserRequest) (string, error) {
	data, err := csvReadAll(l.userFileName)
	if err != nil {
		return "", err
	}

	for _, line := range data {
		if line[0] == req.Name {
			hashed := base64.StdEncoding.EncodeToString([]byte(req.Password))
			if hashed == line[1] {
				return token.Generate() // Generate JWT Token
			}
			logger.Info("wrong password for user: %s", req.Name)
			return "", ErrAuthFailed
		}
	}

	logger.Info("no such user %s", req.Name)
	return "", ErrAuthFailed
}

func (l *localDBHandler) CreateUser(newUser UserRequest) error {
	// User is already exists?
	data, err := csvReadAll(l.userFileName)
	if err != nil {
		return err
	}

	for _, line := range data {
		if line[0] == newUser.Name {
			return ErrUserAlreadyExists
		}
	}

	// add new user
	file, err := os.OpenFile(l.userFileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("Failed to open file %s for append new user", l.userFileName)
		return err
	}
	defer file.Close()

	l.mu.Lock()
	hashed := base64.StdEncoding.EncodeToString([]byte(newUser.Password))
	fmt.Fprintf(file, "%s,%s", newUser.Name, hashed)
	l.mu.Unlock()

	logger.Info("User %s is successfully created", newUser.Name)
	return nil
}

func (l *localDBHandler) DeleteUser(userName string) error {
	var data [][]string

	file, err := os.OpenFile(l.userFileName, os.O_RDWR, 0644)
	if err != nil {
		logger.Error("Failed to open DB file %s in Delete: %v", l.userFileName, err)
		return err
	}
	defer file.Close()

	isDeleted := false
	reader := csv.NewReader(file)

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("Failed to read data: %v", err)
			return err
		}
		if line[0] == userName {
			isDeleted = true
		} else {
			data = append(data, line)
		}
	}

	l.mu.Lock()

	// Remove All data at first
	file.Truncate(0)
	file.Seek(0, 0)

	writer := csv.NewWriter(file)
	writer.WriteAll(data)

	l.mu.Unlock()

	if !isDeleted {
		logger.Info("no such user %s", userName)
		return ErrNoSuchUser
	}

	logger.Info("User %s is successfully delete", userName)
	return nil
}
