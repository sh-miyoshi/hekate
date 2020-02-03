package userapi

// UserCreateRequest ...
type UserCreateRequest struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	SystemRoles []string `json:"system_roles"`
	CustomRoles []string `json:"custom_roles"`
}

// UserGetResponse ...
type UserGetResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	PasswordHash string   `json:"passwordHash"`
	CreatedAt    string   `json:"createdAt"`
	SystemRoles  []string `json:"system_roles"`
	CustomRoles  []string `json:"custom_roles"`
	Sessions     []string `json:"sessions"` // Array of session IDs
}

// UserPutRequest ...
type UserPutRequest struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	SystemRoles []string `json:"system_roles"`
	CustomRoles []string `json:"custom_roles"`
}
