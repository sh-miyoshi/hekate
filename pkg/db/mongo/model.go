package mongo

import (
	"time"
)

type tokenConfig struct {
	AccessTokenLifeSpan  uint   `bson:"accessTokenLifeSpan"`
	RefreshTokenLifeSpan uint   `bson:"refreshTokenLifeSpan"`
	SigningAlgorithm     string `bson:"signingAlgorithm"`
	SignPublicKey        []byte `bson:"signPublicKey"`
	SignSecretKey        []byte `bson:"signSecretKey"`
}

type passwordPolicy struct {
	MinimumLength       uint     `bson:"length"`
	NotUserName         bool     `bson:"notUserName"`
	BlackList           []string `bson:"blackList"`
	UseCharacter        string   `bson:"useCharacter"`
	UseDigit            bool     `bson:"useDigit"`
	UseSpecialCharacter bool     `bson:"useSpecialCharacter"`
}

type projectInfo struct {
	Name            string         `bson:"name"`
	CreatedAt       time.Time      `bson:"createAt"`
	TokenConfig     *tokenConfig   `bson:"tokenConfig"`
	PermitDelete    bool           `bson:"permitDelete"`
	AllowGrantTypes []string       `bson:"allowGrantTypes"`
	PasswordPolicy  passwordPolicy `bson:"passwordPolicy"`
}

type session struct {
	UserID       string    `bson:"userID"`
	ProjectName  string    `bson:"projectName"`
	SessionID    string    `bson:"sessionID"`
	CreatedAt    time.Time `bson:"createdAt"`
	ExpiresIn    uint      `bson:"expiresIn"`
	FromIP       string    `bson:"fromIP"`
	LastAuthTime time.Time `bson:"lastAuthTime"`
	AuthMaxAge   uint      `bson:"maxAge"`
}

type authCodeSession struct {
	SessionID    string    `bson:"sessionID"`
	Code         string    `bson:"code"`
	ExpiresIn    time.Time `bson:"expiresIn"`
	Scope        string    `bson:"scope"`
	ResponseType []string  `bson:"responseType"`
	ClientID     string    `bson:"clientID"`
	RedirectURI  string    `bson:"redirectURI"`
	Nonce        string    `bson:"nonce"`
	ProjectName  string    `bson:"projectName"`
	MaxAge       uint      `bson:"maxAge"`
	ResponseMode string    `bson:"responseMode"`
	Prompt       []string  `bson:"prompt"`
	LoginDate    time.Time `bson:"loginDate"`
}

type userInfo struct {
	ID           string    `bson:"id"`
	ProjectName  string    `bson:"projectName"`
	Name         string    `bson:"name"`
	CreatedAt    time.Time `bson:"createdAt"`
	PasswordHash string    `bson:"passwordHash"`
	SystemRoles  []string  `bson:"systemRoles"`
	CustomRoles  []string  `bson:"customRoles"`
}

type clientInfo struct {
	ID                  string    `bson:"id"`
	ProjectName         string    `bson:"projectName"`
	Secret              string    `bson:"secret"`
	AccessType          string    `bson:"accessType"`
	CreatedAt           time.Time `bson:"createdAt"`
	AllowedCallbackURLs []string  `bson:"allowedCallbackURLs"`
}

type authCode struct {
	CodeID      string    `bson:"codeID"`
	ExpiresIn   time.Time `bson:"expiresIn"`
	ClientID    string    `bson:"clientID"`
	RedirectURL string    `bson:"redirectURL"`
	UserID      string    `bson:"userID"`
	Nonce       string    `bson:"nonce"`
	MaxAge      uint      `bson:"maxAge"`
}

type customRole struct {
	ID          string    `bson:"id"`
	Name        string    `bson:"name"`
	CreatedAt   time.Time `bson:"createdAt"`
	ProjectName string    `bson:"projectName"`
}

type customRoleInUser struct {
	ProjectName  string `bson:"projectName"`
	UserID       string `bson:"userID"`
	CustomRoleID string `bson:"customRoleID"`
}
