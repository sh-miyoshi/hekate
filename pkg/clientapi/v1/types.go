package clientapi

// ClientCreateRequest ...
type ClientCreateRequest struct {
	ID         string `json:"id"`
	Secret     string `json:"secret"`
	AccessType string `json:"access_type"`
}

// ClientGetResponse ...
type ClientGetResponse struct {
	ID         string `json:"id"`
	Secret     string `json:"secret"`
	AccessType string `json:"access_type"`
	CreatedAt  string `json:"createdAt"`
}

// ClientPutRequest ...
type ClientPutRequest struct {
	Secret     string `json:"secret"`
	AccessType string `json:"access_type"`
}
