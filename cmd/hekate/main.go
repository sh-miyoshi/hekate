package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

func main() {
	// initialize config
	if err := config.InitConfig(os.Args); err != nil {
		errors.Print(errors.Append(err, "Failed to init config"))
		os.Exit(1)
	}

	// Initialize server
	if err := initAll(); err != nil {
		errors.Print(errors.Append(err, "Failed to init server"))
		os.Exit(1)
	}

	// Setup API
	r := mux.NewRouter()
	setAPI(r)

	// Run Database GC
	go db.RunGC()

	cfg := config.Get()

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
