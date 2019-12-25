package login

import (
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
)

// Do ...
func Do(serverAddr string, projectName string, req *tokenapi.TokenRequest) (*tokenapi.TokenResponse,error) {
	url := fmt.Sprintf("%s/api/v1/project/%s/token", serverAddr, projectName)
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to request server: %v", err)
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 200:
		var res tokenapi.TokenResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, fmt.Errorf("<Program Bug> Failed to parse http response: %v", err)
		}

		return &res, nil
	case 401, 404:
		return nil, fmt.Errorf("Failed to login system\nPlease cheak user name or password (or project name)")
	case 500:
		return nil, fmt.Errorf("Internal Server Error is occured\nPlease contact to your server administrator")
	default:
		return nil, fmt.Errorf("<Program Bug> Unexpected http response code: %d", httpRes.StatusCode)
	}
}