package userapi

import (
	"time"
)

// UserCreateRequest ...
type UserCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

// UserGetResponse ...
type UserGetResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Enabled      bool      `json:"enabled"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
	Roles        []string  `json:"roles"` // Array of role IDs
}

// UserPutRequest ...
type UserPutRequest struct {
	Name     string   `json:"name"`
	Enabled  bool     `json:"enabled"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"` // Array of role IDs
}
