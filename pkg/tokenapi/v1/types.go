package tokenapi

// TokenCreateRequest is a struct for create request of token
type TokenCreateRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// TokenCreateResponse is a struct for response of token create
type TokenCreateResponse struct {
	Token string `json:"token"`
}
