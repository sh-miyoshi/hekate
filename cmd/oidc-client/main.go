package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

type secretInfo struct {
	Issuer       string `yaml:"issuer"`
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	RedirectURL  string `yaml:"redirect-url"`
}

var secret secretInfo

func setAPI(r *mux.Router) {
	// Main API ( require auth )
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This is test code, so check bearer token in cookie
		if _, err := r.Cookie("Authorization"); err != nil {
			config, _ := getOIDCConfig()
			url := config.AuthCodeURL("", oidc.Nonce("jhgrgw3iohgor4jioh"))
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		w.Write([]byte("login success"))
	}).Methods("GET")

	r.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		config, provider := getOIDCConfig()
		if err := r.ParseForm(); err != nil {
			http.Error(w, "parse form error", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Request: %v\n", r)

		accessToken, err := config.Exchange(context.Background(), r.Form.Get("code"))
		if err != nil {
			msg := fmt.Sprintf("Can't get access token: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		fmt.Printf("Access Token: %v\n", accessToken)

		rawIDToken, ok := accessToken.Extra("id_token").(string)
		fmt.Printf("ID Token: %s\n", rawIDToken)

		if !ok {
			http.Error(w, "missing token", http.StatusInternalServerError)
			return
		}
		oidcConfig := &oidc.Config{
			ClientID: secret.ClientID,
		}
		verifier := provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			http.Error(w, "id token verify error", http.StatusInternalServerError)
			return
		}
		// dump all id token
		idTokenClaims := map[string]interface{}{}
		if err := idToken.Claims(&idTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("%#v\n", idTokenClaims)
		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: "Bearer " + rawIDToken,
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods("GET")
}

func initSecret(filePath string) error {
	fp, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(&secret); err != nil {
		return err
	}

	return nil
}

func getOIDCConfig() (*oauth2.Config, *oidc.Provider) {
	provider, err := oidc.NewProvider(context.Background(), secret.Issuer)
	if err != nil {
		panic(err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     secret.ClientID,
		ClientSecret: secret.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
		RedirectURL:  fmt.Sprintf("http://%s:3000/callback", secret.RedirectURL),
	}

	return oauth2Config, provider
}

func main() {
	// read secret file
	if err := initSecret("secret.yaml"); err != nil {
		fmt.Printf("Failed to initialize secret: %v", err)
		return
	}

	// Setup API
	r := mux.NewRouter()
	setAPI(r)

	// Run Server
	corsObj := handlers.AllowedOrigins([]string{"*"})
	addr := "0.0.0.0:3000"
	fmt.Printf("start server with %s\n", addr)
	if err := http.ListenAndServe(addr, handlers.CORS(corsObj)(r)); err != nil {
		fmt.Printf("Failed to run server: %+v\n", err)
		os.Exit(1)
	}
}
