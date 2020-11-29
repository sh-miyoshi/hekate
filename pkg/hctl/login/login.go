package login

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
)

func createClient(serverName string, insecure bool, timeout time.Duration) *http.Client {
	print.Debug("request server name: %s, insecure: %v, timeout: %v", serverName, insecure, timeout)
	tlsConfig := tls.Config{
		ServerName: serverName,
	}

	if insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return client
}

func tokenRequest(req *http.Request, insecure bool, timeout uint) (*oidcapi.TokenResponse, error) {
	print.Debug("token request to %s", req.URL.String())

	client := createClient(req.Host, insecure, time.Duration(timeout)*time.Second)
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
		b, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("<Program Bug> Unexpected http response code: %d, message: %s", httpRes.StatusCode, b)
	}
}

// Do ...
func Do(userName, password string, info Info) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/openid-connect/token", info.ServerAddr, info.ProjectName)

	form := url.Values{}
	form.Add("username", userName)
	form.Add("password", password)
	form.Add("grant_type", "password")
	form.Add("client_id", info.ClientID)
	if info.ClientSecret != "" {
		form.Add("client_secret", info.ClientSecret)
	}

	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tokenRequest(httpReq, info.Insecure, info.Timeout)
}

// DoWithRefresh ...
func DoWithRefresh(refreshToken string, info Info) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/api/v1/project/%s/openid-connect/token", info.ServerAddr, info.ProjectName)

	form := url.Values{}
	form.Add("refresh_token", refreshToken)
	form.Add("grant_type", "refresh_token")
	form.Add("client_id", info.ClientID)
	if info.ClientSecret != "" {
		form.Add("client_secret", info.ClientSecret)
	}

	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tokenRequest(httpReq, info.Insecure, info.Timeout)
}
