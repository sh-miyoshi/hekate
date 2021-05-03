package customroleapi

// CustomRoleCreateRequest ...
type CustomRoleCreateRequest struct {
	Name string `json:"name"`
}

// CustomRoleGetResponse ...
type CustomRoleGetResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	CreatedAt   string `json:"created_at"`
}

// CustomRolePutRequest ...
type CustomRolePutRequest struct {
	Name string `json:"name"`
}
