package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {
	serverAddr := "http://localhost:18443"
	issuer := serverAddr + "/api/v1/project/master"
	clientID := "portal"
	clientSecret := ""
	state := "mystate"

	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		fmt.Printf("Failed to create provider: %v\n", err)
		os.Exit(1)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
		RedirectURL:  "http://localhost:3000/callback",
	}

	fmt.Printf("config: %v\n", config)

	values := neturl.Values{}
	values.Set("scope", "openid")
	values.Set("response_type", "id_token")
	values.Set("client_id", clientID)
	values.Set("redirect_uri", config.RedirectURL)
	values.Set("state", state)
	// TODO: values.Set("nonce",""), values.Set("max_age",""), values.Set("prompt","")
	req, err := http.NewRequest("GET", config.Endpoint.AuthURL, nil)
	if err != nil {
		fmt.Printf("Failed to create new request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to server: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// maybe return user login page
	fmt.Printf("auth res: %v\n", res)

	bytes, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`/*action="[^\"]+`)
	url := re.FindString(string(bytes))
	url = strings.TrimPrefix(url, "action=")
	url = strings.Trim(url, "\"")
	url = serverAddr + url

	// fmt.Printf("login user page: %s\n", string(bytes))

	fmt.Printf("redirect url: %s\n", url)
	u, _ := neturl.Parse(url)
	code := u.Query().Get("login_session_id")
	fmt.Printf("code: %s\n", code)

	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	values = neturl.Values{
		"username":         []string{"admin"},
		"password":         []string{"password"},
		"login_session_id": []string{code},
	}
	req.URL.RawQuery = values.Encode()

	client = &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err = client.Do(req)
	if err != nil {
		fmt.Printf("Failed to request to user login handler: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	fmt.Printf("login res: %v\n", res)

	if res.StatusCode != http.StatusFound {
		fmt.Printf("Unexpected login result: want status %d, got %d\n", http.StatusFound, res.StatusCode)
		os.Exit(1)
	}

	// get id token from res
	url = res.Header.Get("Location")
	if !strings.Contains(url, "id_token=") {
		fmt.Println("Failed to get id token")
		os.Exit(1)
	}
	fmt.Println("Success to get id token")
}
