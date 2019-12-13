package userapi

// UserCreateRequest ...
type UserCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

// UserGetResponse ...
type UserGetResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	PasswordHash string   `json:"passwordHash"`
	CreatedAt    string   `json:"createdAt"`
	Roles        []string `json:"roles"`    // Array of role IDs
	Sessions     []string `json:"sessions"` // Array of session IDs
}

// UserPutRequest ...
type UserPutRequest struct {
	Name     string   `json:"name"`
	Enabled  bool     `json:"enabled"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"` // Array of role IDs
}
