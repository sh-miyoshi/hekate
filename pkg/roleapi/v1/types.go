package roleapi

// RoleCreateRequest ...
type RoleCreateRequest struct {
	Name           string   `json:"name"`
	TargetResource []string `json:"targetResource"`
	Type           string   `json:"type"`
}

// RoleGetResponse ...
type RoleGetResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	TargetResource []string `json:"targetResource"`
	Type           string   `json:"type"`
	CreatedAt      string   `json:"createdAt"`
}

// RolePutRequest ...
type RolePutRequest struct {
	Name           string   `json:"name,omitempty"`
	TargetResource []string `json:"targetResource,omitempty"`
	Type           string   `json:"type,omitempty"`
}
