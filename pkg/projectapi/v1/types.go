package projectapi

import (
	"time"
)

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
	Name        string       `json:"name"`
	Enabled     bool         `json:"enabled"`
	CreatedAt   time.Time    `json:"createdAt"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}

// ProjectPutRequest ...
type ProjectPutRequest struct {
	Enabled     bool         `json:"enabled"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}
