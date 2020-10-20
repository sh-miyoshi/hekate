package sessionapi

// SessionGetResponse ...
type SessionGetResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	ExpiresIn int64  `json:"expires_in"`
	FromIP    string `json:"from_ip"`
}
