package tokenapi

// TokenRequest ...
//   require: name or id
type TokenRequest struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Secret   string `json:"secret"`
	AuthType string `json:"authType"`
}

// TokenResponse ...
type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  uint   `json:"accessExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn uint   `json:"refreshExpiresIn"`
}
