package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
)

func tokenRequest(req *http.Request) (*oidcapi.TokenResponse, error) {
	client := &http.Client{}
	httpRes, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to request server: %v", err)
	}
	defer httpRes.Body.Close()

	switch httpRes.StatusCode {
	case 200:
		var res oidcapi.TokenResponse
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

// Do ...
func Do(serverAddr string, projectName string, userName string, password string) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/openid-connect/token", serverAddr, projectName)

	form := url.Values{}
	form.Add("username", userName)
	form.Add("password", password)
	form.Add("grant_type", "password")
	form.Add("client_id", "admin-cli")
	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tokenRequest(httpReq)
}

// DoWithRefresh ...
func DoWithRefresh(serverAddr string, projectName string, refreshToken string) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/openid-connect/token", serverAddr, projectName)

	form := url.Values{}
	form.Add("refresh_token", refreshToken)
	form.Add("grant_type", "refresh_token")
	form.Add("client_id", "admin-cli")
	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tokenRequest(httpReq)
}
