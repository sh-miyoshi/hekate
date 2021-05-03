package keysapi

// KeysGetResponse ...
type KeysGetResponse struct {
	Type      string `json:"type"`
	PublicKey string `json:"public_key"`
}
