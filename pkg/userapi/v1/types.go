package userapi

// UserCreateRequest ...
type UserCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UserGetResponse ...
type UserGetResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PasswordHash string   `json:"passwordHash"`
	CreatedAt    string   `json:"createdAt"`
	Roles        []string `json:"roles"`    // Array of role IDs
	Sessions     []string `json:"sessions"` // Array of session IDs
}

// UserPutRequest ...
type UserPutRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"` // Array of role IDs
}
