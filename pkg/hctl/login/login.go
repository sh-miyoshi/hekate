package login

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
	"strings"
	"time"

	oauthapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/oauth"
	"github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/oidc"
	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/oidc"
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

func reqHandler(method string, url string, form neturl.Values, insecure bool, timeout uint, resHandler func(httpRes *http.Response) (interface{}, error)) (interface{}, error) {
	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := createClient(httpReq.Host, insecure, time.Duration(timeout)*time.Second)
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to request server: %v", err)
	}

	dump, _ := httputil.DumpResponse(httpRes, true)
	print.Debug("Response: %q", dump)

	return resHandler(httpRes)
}

// Do method login by device flow
//   start login session
//   show device login page and user code
//   polling to token endpoint until user login by web browser
func Do(info Info) (*oidcapi.TokenResponse, error) {
	url := fmt.Sprintf("%s/adminapi/v1/project/%s/oauth/device", info.ServerAddr, info.ProjectName)
	form := neturl.Values{}
	form.Add("scope", "openid")
	form.Add("client_id", info.ClientID)
	if info.ClientSecret != "" {
		form.Add("client_secret", info.ClientSecret)
	}
	rawRes, err := reqHandler("POST", url, form, info.Insecure, info.Timeout, func(httpRes *http.Response) (interface{}, error) {
		defer httpRes.Body.Close()

		switch httpRes.StatusCode {
		case 200:
			var res oauthapi.DeviceAuthorizationResponse
			// parse response
			if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("<Program Bug> Failed to parse http response: %v", err)
			}
			return res, nil
		case 400:
			b, _ := ioutil.ReadAll(httpRes.Body)
			return nil, fmt.Errorf("Failed to login system: %s", b)
		case 500:
			return nil, fmt.Errorf("Internal Server Error is occured\nPlease contact to your server administrator")
		}
		b, _ := ioutil.ReadAll(httpRes.Body)
		return nil, fmt.Errorf("<Program Bug> Unexpected http response code: %d, message: %s", httpRes.StatusCode, b)
	})
	if err != nil {
		return nil, err
	}
	res := rawRes.(oauthapi.DeviceAuthorizationResponse)

	print.Print("To log in, use a web browser to open the page %s and enter the code %s to authenticate.", res.VerificationURI, res.UserCode)

	interval := 5 // 5[sec]
	if res.Interval != 0 {
		interval = res.Interval
	}

	url = fmt.Sprintf("%s/adminapi/v1/project/%s/openid-connect/token", info.ServerAddr, info.ProjectName)
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Add("device_code", res.DeviceCode)

	// Get token
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		rawRes, err := reqHandler("POST", url, form, info.Insecure, info.Timeout, func(httpRes *http.Response) (interface{}, error) {
			defer httpRes.Body.Close()

			switch httpRes.StatusCode {
			case 200:
				var res oidcapi.TokenResponse
				if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
					return nil, fmt.Errorf("<Program Bug> Failed to parse http response: %v", err)
				}
				return res, nil
			case 400:
				var res oidc.ErrorResponse
				if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
					return nil, fmt.Errorf("<Program Bug> Failed to parse http response: %v", err)
				}

				// waiting authorization
				if res.ErrorCode == "authorization_pending" {
					return nil, nil
				}
				if res.ErrorCode == "slow_down" {
					interval++
					return nil, nil
				}
			case 500:
				return nil, fmt.Errorf("Internal Server Error is occured\nPlease contact to your server administrator")
			}
			b, _ := ioutil.ReadAll(httpRes.Body)
			return nil, fmt.Errorf("<Program Bug> Unexpected http response code: %d, message: %s", httpRes.StatusCode, b)
		})
		if err != nil {
			return nil, err
		}

		if rawRes == nil {
			continue
		}

		res := rawRes.(oidcapi.TokenResponse)
		return &res, nil
	}
}

// DoWithRefresh ...
func DoWithRefresh(refreshToken string, info Info) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/adminapi/v1/project/%s/openid-connect/token", info.ServerAddr, info.ProjectName)

	form := neturl.Values{}
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

// DoWithClient ...
func DoWithClient(info Info) (*oidcapi.TokenResponse, error) {
	u := fmt.Sprintf("%s/adminapi/v1/project/%s/openid-connect/token", info.ServerAddr, info.ProjectName)

	form := neturl.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", info.ClientID)
	form.Add("client_secret", info.ClientSecret)

	body := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, fmt.Errorf("<Program Bug> Failed to create http request: %v", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return tokenRequest(httpReq, info.Insecure, info.Timeout)
}
