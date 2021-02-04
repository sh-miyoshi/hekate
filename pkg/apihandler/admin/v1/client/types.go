package clientapi

// ClientCreateRequest ...
type ClientCreateRequest struct {
	ID                  string   `json:"id"`
	Secret              string   `json:"secret"`
	AccessType          string   `json:"access_type"`
	AllowedCallbackURLs []string `json:"allowed_callback_urls"`
}

// ClientGetResponse ...
type ClientGetResponse struct {
	ID                  string   `json:"id"`
	Secret              string   `json:"secret"`
	AccessType          string   `json:"access_type"`
	CreatedAt           string   `json:"created_at"`
	AllowedCallbackURLs []string `json:"allowed_callback_urls"`
}

// ClientPutRequest ...
type ClientPutRequest struct {
	Secret              string   `json:"secret"`
	AccessType          string   `json:"access_type"`
	AllowedCallbackURLs []string `json:"allowed_callback_urls"`
}
