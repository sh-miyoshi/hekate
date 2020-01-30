package main

import (
	"context"
	"fmt"
	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"os"
	"net/http"
)

func main() {
	issuer := "http://localhost:8080/api/v1/project/master"
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

	// maybe return user login page
	fmt.Printf("auth res: %v\n", res)


}
