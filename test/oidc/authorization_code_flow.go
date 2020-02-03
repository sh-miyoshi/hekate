package main

import (
	"context"
	"fmt"
	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"os"
	"net/http"
	"regexp"
	"io/ioutil"
	"strings"
	neturl "net/url"
)

func main() {
	serverAddr := "http://localhost:8080"
	issuer := serverAddr + "/api/v1/project/master"
	clientID := "admin-cli"
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

	url := config.AuthCodeURL(state)

	fmt.Printf("url: %v\n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to request to auth handler: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// maybe return user login page
	fmt.Printf("auth res: %v\n", res)

	bytes, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`/*action="[^\"]+`)
	url = re.FindString(string(bytes))
	url = strings.TrimPrefix(url, "action=")
	url = strings.Trim(url, "\"")
	url = serverAddr + url

	// fmt.Printf("login user page: %s\n", string(bytes))

	fmt.Printf("redirect url: %s\n", url)
	u, _ := neturl.Parse(url)
	code := u.Query().Get("login_verify_code")
	fmt.Printf("code: %s\n", code)

	req,_ := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	values := neturl.Values{
		"username": []string{"admin"},
		"password": []string{"password"},
		"login_verify_code": []string{code},
	}
	req.URL.RawQuery = values.Encode()
	
	client := &http.Client{}
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

	// TODO(check status code: 302 -> ok)
	// TODO(exchange code to token)
}
