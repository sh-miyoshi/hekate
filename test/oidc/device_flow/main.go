package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strings"

	oauthapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/auth/v1/oauth"
	"github.com/sh-miyoshi/hekate/pkg/apihandler/auth/v1/oidc"
)

func main() {
	serverAddr := "http://localhost:18443"
	clientID := "portal"

	values := neturl.Values{}
	values.Set("scope", "openid")
	values.Set("client_id", clientID)
	url := serverAddr + "/adminapi/v1/project/master/oauth/device"
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

	// send user code
	values = neturl.Values{}
	values.Set("code", res.UserCode)
	url = serverAddr + "/resource/project/master/deviceverify"
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Failed to create new request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()
	httpRes, err = client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to server: %v\n", err)
		os.Exit(1)
	}

	// maybe return login page
	if httpRes.StatusCode != 200 {
		fmt.Printf("Failed to get login page: %s\n", httpRes.Status)
		os.Exit(1)
	}

	// login by web browser
	bytes, _ := ioutil.ReadAll(httpRes.Body)
	httpRes.Body.Close()
	re := regexp.MustCompile(`/*action="[^\"]+`)
	url = re.FindString(string(bytes))
	url = strings.TrimPrefix(url, "action=")
	url = strings.Trim(url, "\"")
	url = serverAddr + url

	// fmt.Printf("login user page: %s\n", string(bytes))

	fmt.Printf("redirect url: %s\n", url)
	u, _ := neturl.Parse(url)
	code := u.Query().Get("login_session_id")
	fmt.Printf("login session id: %s\n", code)

	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	values = neturl.Values{
		"username":         []string{"admin"},
		"password":         []string{"password"},
		"login_session_id": []string{code},
	}
	req.URL.RawQuery = values.Encode()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	httpRes, err = client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to server: %v\n", err)
		os.Exit(1)
	}
	if httpRes.StatusCode != http.StatusFound {
		fmt.Printf("Unexpected login result: want status %d, got %d\n", http.StatusFound, httpRes.StatusCode)
		os.Exit(1)
	}

	values = neturl.Values{}
	values.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	values.Set("client_id", clientID)
	values.Set("device_code", res.DeviceCode)
	url = serverAddr + "/authapi/v1/project/master/openid-connect/token"
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	httpRes, err = client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to server: %v\n", err)
		os.Exit(1)
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var successRes oidc.TokenResponse
		json.NewDecoder(httpRes.Body).Decode(&successRes)
		fmt.Printf("Access Token: %s\n", successRes.AccessToken)
		fmt.Printf("Refresh Token: %s\n", successRes.RefreshToken)
	} else if httpRes.StatusCode == http.StatusBadRequest {
		var errRes oidc.ErrorResponse
		json.NewDecoder(httpRes.Body).Decode(&errRes)
		fmt.Printf("Error message: %s\n", errRes.ErrorCode)
		os.Exit(1)
	} else {
		fmt.Printf("Unexpected error: %s\n", httpRes.Status)
		os.Exit(1)
	}
}
