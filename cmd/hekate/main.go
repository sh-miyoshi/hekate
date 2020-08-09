package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sh-miyoshi/hekate/cmd/hekate/config"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

const (
	authCodeUserLoginResourcePath = "/resource/login"
)

func main() {
	// Get config
	cfg, err := config.InitConfig(os.Args)
	if err != nil {
		fmt.Printf("Failed to set config: %v", err)
		os.Exit(1)
	}

	// Initialize server
	if err := initAll(cfg); err != nil {
		errors.Print(errors.Append(err, "Failed to init server"))
		os.Exit(1)
	}

	// Setup API
	r := mux.NewRouter()
	setAPI(r, cfg)

	// Run Server
	addr := fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.Port)
	logger.Info("start server with %s", addr)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
	})

	if cfg.HTTPSConfig.Enabled {
		logger.Info("Run server as https")
		if err := http.ListenAndServeTLS(addr, cfg.HTTPSConfig.CertFile, cfg.HTTPSConfig.KeyFile, corsOpts.Handler(r)); err != nil {
			logger.Error("Failed to run server: %v", err)
			os.Exit(1)
		}
	} else {
		logger.Info("Run server as http")
		if err := http.ListenAndServe(addr, corsOpts.Handler(r)); err != nil {
			logger.Error("Failed to run server: %v", err)
			os.Exit(1)
		}
	}
}
