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

type projectInfo struct {
	Name            string       `bson:"name"`
	CreatedAt       time.Time    `bson:"createAt"`
	TokenConfig     *tokenConfig `bson:"tokenConfig"`
	PermitDelete    bool         `bson:"permitDelete"`
	AllowGrantTypes []string     `bson:"allowGrantTypes"`
}

type session struct {
	UserID    string    `bson:"userID"`
	SessionID string    `bson:"sessionID"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresIn uint      `bson:"expiresIn"`
	FromIP    string    `bson:"fromIP"`
}

type loginSessionInfo struct {
	VerifyCode   string    `bson:"verifyCode"`
	ExpiresIn    time.Time `bson:"expiresIn"`
	Scope        string    `bson:"scope"`
	ResponseType string    `bson:"responseType"`
	ClientID     string    `bson:"clientID"`
	RedirectURI  string    `bson:"redirectURI"`
	Nonce        string    `bson:"nonce"`
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
}

type customRole struct {
	ID          string    `bson:"id"`
	Name        string    `bson:"name"`
	CreatedAt   time.Time `bson:"createdAt"`
	ProjectName string    `bson:"projectName"`
}
