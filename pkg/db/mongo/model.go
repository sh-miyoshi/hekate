package mongo

import (
	"time"
)

type tokenConfig struct {
	AccessTokenLifeSpan  uint `bson:"accessTokenLifeSpan"`
	RefreshTokenLifeSpan uint `bson:"refreshTokenLifeSpan"`
}

type projectInfo struct {
	Name        string       `bson:"name"`
	CreatedAt   time.Time    `bson:"createAt"`
	TokenConfig *tokenConfig `bson:"tokenConfig"`
}
