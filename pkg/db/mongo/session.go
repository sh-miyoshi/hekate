package mongo

import (
	// "github.com/pkg/errors"
	// "github.com/sh-miyoshi/jwt-server/pkg/db/model"
	// "time"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SessionHandler implement db.SessionHandler
type SessionHandler struct {
	dbClient *mongo.Client
}

// NewSessionHandler ...
func NewSessionHandler(dbClient *mongo.Client) *SessionHandler {
	res := &SessionHandler{
		dbClient: dbClient,
	}
	return res
}

// New ...
func (h *SessionHandler) New(userID string, sessionID string, expiresIn uint, fromIP string) error {
	return nil
}

// Revoke ...
func (h *SessionHandler) Revoke(sessionID string) error {
	return nil
}

// GetList ...
func (h *SessionHandler) GetList(userID string) ([]string, error) {
	res := []string{}

	return res, nil
}
