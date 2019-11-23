package projectapi

// TODO(move to other directory)

// TokenConfig ...
type TokenConfig struct {
	AccessTokenLifeSpan  int32 `json:"accessTokenLifeSpan"`
	RefreshTokenLifeSpan int32 `json:"refreshTokenLifeSpan"`
}

// ProjectCreateRequest ...
type ProjectCreateRequest struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled,omitempty"`
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
	Name        string       `json:"name,omitempty"`
	Enabled     bool         `json:"enabled,omitempty"`
	TokenConfig *TokenConfig `json:"tokenConfig,omitempty"`
}
