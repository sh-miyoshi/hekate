package db

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

const (
	roleSeparator string = ";"
)

type localDBHandler struct {
	DBHandler

	dbFileName string
	nextID     int
	mu         sync.Mutex
}

// This func read all csv data at once, so should not use in production
func (l *localDBHandler) csvReadAll() ([][]string, error) {
	file, err := os.Open(l.dbFileName)
	if err != nil {
		logger.Error("Failed to open DB file %s in Authenticate: %v", l.dbFileName, err)
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'

	return reader.ReadAll()
}

func (l *localDBHandler) saveUsers(users []User) error {
	// Cleanup exists data
	file, err := os.Create(l.dbFileName)
	if err != nil {
		return err
	}

	for _, user := range users {
		roles := ""
		for _, role := range user.Roles {
			roles += fmt.Sprintf("%d%s", role, roleSeparator)
		}
		roles = strings.TrimSuffix(roles, roleSeparator)
		fmt.Fprintf(file, "%s,%s,%s,%s", user.ID, user.Name, user.Password, roles)
	}

	return nil
}

// ConnectDB check file exists(connectString is a file path of user data)
func (l *localDBHandler) ConnectDB(connectString string) error {

	// create new db file
	if _, err := os.Create(connectString); err != nil {
		return err
	}

	// TODO check file broken

	l.dbFileName = connectString
	l.nextID = 0
	return nil
}

func (l *localDBHandler) CreateUser(newUser UserRequest) error {
	// User is already exists?
	users, err := l.GetUserList()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == newUser.Name {
			return ErrUserAlreadyExists
		}
	}

	// add new user
	hashed := base64.StdEncoding.EncodeToString([]byte(newUser.Password))
	users = append(users, User{
		ID:       strconv.Itoa(l.nextID),
		Name:     newUser.Name,
		Password: hashed,
	})

	l.mu.Lock()
	l.saveUsers(users)
	l.nextID++
	l.mu.Unlock()

	logger.Info("User %s is successfully created", newUser.Name)
	return nil
}

func (l *localDBHandler) DeleteUser(userID string) error {
	users, err := l.GetUserList()
	if err != nil {
		return err
	}

	isDeleted := false
	newUsers := []User{}

	for _, user := range users {
		if user.ID == userID {
			isDeleted = true
		} else {
			newUsers = append(newUsers, user)
		}
	}

	l.mu.Lock()
	l.saveUsers(users)
	l.mu.Unlock()

	if !isDeleted {
		logger.Info("no such user %s", userID)
		return ErrNoSuchUser
	}

	logger.Info("User %s is successfully delete", userID)
	return nil
}

func (l *localDBHandler) GetUserList() ([]User, error) {
	ret := []User{}

	data, err := l.csvReadAll()
	if err != nil {
		return ret, err
	}

	for _, line := range data {
		roles := []RoleType{}
		if len(line) == 3 { // If data have roles
			roleStrs := strings.Split(line[3], roleSeparator)
			for _, roleStr := range roleStrs {
				role, err := strconv.Atoi(roleStr)
				if err != nil {
					return []User{}, err
				}
				roles = append(roles, RoleType(role))
			}
		}
		ret = append(ret, User{
			ID:       line[0],
			Name:     line[1],
			Password: line[2],
			Roles:    roles,
		})
	}
	return ret, nil
}

func (l *localDBHandler) UpdatePassowrd(newPassword string) error {
	// Not Implemented yet
	return nil
}

func (l *localDBHandler) AddRoleToUser(addRole RoleType, userID string) error {
	users, err := l.GetUserList()
	if err != nil {
		return err
	}
	resUsers := []User{}

	findUser := false
	for _, user := range users {
		if user.ID == userID {
			for _, role := range user.Roles {
				if role == addRole {
					return fmt.Errorf("Adding role is already exists")
				}
			}
			user.Roles = append(user.Roles, addRole)
			findUser = true
		}
		resUsers = append(resUsers, user)
	}

	l.mu.Lock()
	l.saveUsers(resUsers)
	l.mu.Unlock()

	if !findUser {
		return ErrNoSuchUser
	}
	return nil
}
func (l *localDBHandler) RemoveRoleFromUser(role RoleType, userID string) error {
	// Not Implemented yet
	return nil
}

func (l *localDBHandler) SetTokenConfig(config TokenConfig) error {
	// Not Implemented yet
	return nil
}
func (l *localDBHandler) GetTokenConfig() (TokenConfig, error) {
	// Not Implemented yet
	return TokenConfig{}, nil
}
