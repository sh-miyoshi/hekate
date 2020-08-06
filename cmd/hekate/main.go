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
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	defaultrole "github.com/sh-miyoshi/hekate/pkg/role"
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

	// Initialize logger
	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %v", *cfg)

	// Check login resource directory struct
	if err := cfg.CheckLoginResDirStruct(); err != nil {
		errors.Print(errors.Append(err, "Login resource directory is broken"))
		os.Exit(1)
	}

	// Initialize Default Role Handler
	if err := defaultrole.InitHandler(); err != nil {
		errors.Print(errors.Append(err, "Failed to initialize default role handler"))
		os.Exit(1)
	}

	// Initialize Token Config
	token.InitConfig(cfg.HTTPSConfig.Enabled)

	// Initialize OIDC Config
	oidc.InitConfig(cfg.HTTPSConfig.Enabled, cfg.AuthCodeExpiresTime, cfg.UserLoginResourceDir, authCodeUserLoginResourcePath)

	// Initalize Database
	if err := initDB(cfg.DB.Type, cfg.DB.ConnectionString, cfg.AdminName, cfg.AdminPassword); err != nil {
		errors.Print(errors.Append(err, "Failed to initialize database"))
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
