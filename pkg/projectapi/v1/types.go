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
	TokenConfig *TokenConfig `json:"tokenConfig"`
}

// ProjectGetResponse ...
type ProjectGetResponse struct {
	Name        string       `json:"name"`
	CreatedAt   time.Time    `json:"createdAt"`
	TokenConfig *TokenConfig `json:"tokenConfig"`
}

// ProjectPutRequest ...
type ProjectPutRequest struct {
	TokenConfig *TokenConfig `json:"tokenConfig"`
}
