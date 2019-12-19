package mongo

import (
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserInfoHandler implement db.UserInfoHandler
type UserInfoHandler struct {
	dbClient       *mongo.Client
	projectHandler *ProjectInfoHandler
}

// NewUserHandler ...
func NewUserHandler(dbClient *mongo.Client, projectHandler *ProjectInfoHandler) *UserInfoHandler {
	res := &UserInfoHandler{
		projectHandler: projectHandler,
		dbClient:       dbClient,
	}
	return res
}

// Add ...
func (h *UserInfoHandler) Add(ent *model.UserInfo) error {
	return nil
}

// Delete ...
func (h *UserInfoHandler) Delete(projectName string, userID string) error {
	return nil
}

// GetList ...
func (h *UserInfoHandler) GetList(projectName string) ([]string, error) {
	return []string{}, nil
}

// Get ...
func (h *UserInfoHandler) Get(projectName string, userID string) (*model.UserInfo, error) {
	return nil, nil
}

// Update ...
func (h *UserInfoHandler) Update(ent *model.UserInfo) error {
	return nil
}

// GetIDByName ...
func (h *UserInfoHandler) GetIDByName(projectName string, userName string) (string, error) {
	return "", nil
}

// DeleteAll ...
func (h *UserInfoHandler) DeleteAll(projectName string) error {
	return nil
}

// AddRole ...
func (h *UserInfoHandler) AddRole(projectName string, userID string, roleID string) error {
	return nil
}

// DeleteRole ...
func (h *UserInfoHandler) DeleteRole(projectName string, userID string, roleID string) error {
	return nil
}

// NewSession ...
func (h *UserInfoHandler) NewSession(projectName string, userID string, session model.Session) error {
	return nil
}

// RevokeSession ...
func (h *UserInfoHandler) RevokeSession(projectName string, userID string, sessionID string) error {
	return nil
}
