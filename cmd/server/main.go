package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/cmd/server/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"net/http"
	"os"
	"path/filepath"
)

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// Add API
	r.HandleFunc(basePath+"/token", tokenapiv1.TokenCreateHandler).Methods("POST")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func main() {
	const defaultConfigFilePath = "./config.yaml"
	configFilePath := flag.String("config", defaultConfigFilePath, "file name of config.yaml")
	flag.Parse()

	configAbsFilePath, _ := filepath.Abs(*configFilePath)
	cfg, err := config.InitConfig(configAbsFilePath)
	if err != nil {
		fmt.Printf("Failed to set config: %v", err)
		os.Exit(1)
	}

	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %v", *cfg)

	r := mux.NewRouter()
	setAPI(r)

	corsObj := handlers.AllowedOrigins([]string{"*"})

	addr := fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.Port)
	logger.Info("start server with %s", addr)
	if err := http.ListenAndServe(addr, handlers.CORS(corsObj)(r)); err != nil {
		logger.Error("Failed to run server: %v", err)
		os.Exit(1)
	}
}
