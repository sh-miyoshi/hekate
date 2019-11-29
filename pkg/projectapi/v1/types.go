package projectapi

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  int `json:"accessTokenLifeSpan"`
	RefreshTokenLifeSpan int `json:"refreshTokenLifeSpan"`
}

// ProjectCreateRequest ...
type ProjectCreateRequest struct {
	Name        string       `json:"name"`
	Enabled     bool         `json:"enabled"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}

// ProjectGetResponse ...
type ProjectGetResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Enabled     bool         `json:"enabled"`
	CreatedAt   string       `json:"createdAt"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}

// ProjectPutRequest ...
type ProjectPutRequest struct {
	Name        string       `json:"name"`
	Enabled     bool         `json:"enabled"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}
