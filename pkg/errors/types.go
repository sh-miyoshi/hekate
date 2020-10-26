package errors

// Error ...
type Error struct {
	privateInfo      []info
	publicMsg        string
	httpResponseCode int
}

// HTTPError ...
type HTTPError struct {
	Type  string `json:"type"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}
