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

type session struct {
	UserID    string    `bson:"userID"`
	SessionID string    `bson:"sessionID"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresIn uint      `bson:"expiresIn"`
	FromIP    string    `bson:"fromIP"`
}

type userInfo struct {
	ID           string    `bson:"id"`
	ProjectName  string    `bson:"projectName"`
	Name         string    `bson:"name"`
	CreatedAt    time.Time `bson:"createdAt"`
	PasswordHash string    `bson:"passwordHash"`
	Roles        []string  `bson:"roles"`
}

type clientInfo struct {
	ID          string    `bson:"id"`
	ProjectName string    `bson:"projectName"`
	Secret      string    `bson:"secret"`
	AccessType  string    `bson:"accessType"`
	CreatedAt   time.Time `bson:"createdAt"`
}
