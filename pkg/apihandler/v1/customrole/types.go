package customroleapi

// CustomRoleCreateRequest ...
type CustomRoleCreateRequest struct {
	Name string `json:"name"`
}

// CustomRoleGetResponse ...
type CustomRoleGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

// CustomRolePutRequest ...
type CustomRolePutRequest struct {
	Name string `json:"name"`
}
