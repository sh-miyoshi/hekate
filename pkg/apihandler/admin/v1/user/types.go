package userapi

// CustomRole ...
type CustomRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserCreateRequest ...
type UserCreateRequest struct {
	Name        string   `json:"name"`
	Password    string   `json:"password"`
	SystemRoles []string `json:"system_roles"`
	CustomRoles []string `json:"custom_roles"`
}

// UserGetResponse ...
type UserGetResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	CreatedAt   string       `json:"createdAt"`
	SystemRoles []string     `json:"system_roles"`
	CustomRoles []CustomRole `json:"custom_roles"`
	Sessions    []string     `json:"sessions"` // Array of session IDs
	Locked      bool         `json:"locked"`
	// TODO OTP Info
}

// UserPutRequest ...
type UserPutRequest struct {
	Name        string   `json:"name"`
	SystemRoles []string `json:"system_roles"`
	CustomRoles []string `json:"custom_roles"`
}

// UserResetPasswordRequest ...
type UserResetPasswordRequest struct {
	Password string `json:"password"`
}
