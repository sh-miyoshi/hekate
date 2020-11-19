package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"time"

	oauthapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oauth"
	"github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
)

func main() {
	serverAddr := "http://localhost:18443"
	clientID := "portal"

	values := neturl.Values{}
	values.Set("scope", "openid")
	values.Set("client_id", clientID)
	url := serverAddr + "/api/v1/project/master/oauth/device"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Failed to create new request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	client := &http.Client{}
	httpRes, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to server: %v\n", err)
		os.Exit(1)
	}

	res := oauthapiv1.DeviceAuthorizationResponse{}
	json.NewDecoder(httpRes.Body).Decode(&res)
	fmt.Printf("Access URL: %s\n", res.VerificationURI)
	fmt.Printf("User Code: %s\n", res.UserCode)
	httpRes.Body.Close()

	values = neturl.Values{}
	values.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	values.Set("client_id", clientID)
	values.Set("device_code", res.DeviceCode)
	url = serverAddr + "/api/v1/project/master/openid-connect/token"
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	for {
		httpRes, err = client.Do(req)
		if err != nil {
			fmt.Printf("Failed to request to server: %v\n", err)
			os.Exit(1)
		}

		if httpRes.StatusCode == http.StatusOK {
			var successRes oidc.TokenResponse
			json.NewDecoder(httpRes.Body).Decode(&successRes)
			fmt.Printf("Access Token: %s\n", successRes.AccessToken)
			fmt.Printf("Refresh Token: %s\n", successRes.RefreshToken)
			httpRes.Body.Close()
			break
		} else if httpRes.StatusCode == http.StatusBadRequest {
			var errRes oidc.ErrorResponse
			json.NewDecoder(httpRes.Body).Decode(&errRes)
			fmt.Printf("Error message: %s\n", errRes.ErrorCode)
		}
		httpRes.Body.Close()
		time.Sleep(10 * time.Second)
	}
}
