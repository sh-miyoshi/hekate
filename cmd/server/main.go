package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/cmd/server/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	userapiv1 "github.com/sh-miyoshi/jwt-server/pkg/userapi/v1"
	"net/http"
	"os"
	"path/filepath"
)

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// Token API
	r.HandleFunc(basePath+"/token", tokenapiv1.TokenCreateHandler).Methods("POST")

	// Project API

	// User API
	r.HandleFunc(basePath+"/project/{projectID}/user", userapiv1.AllUserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/user", userapiv1.UserCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserUpdateHandler).Methods("PUT")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}/role/{roleID}", userapiv1.UserRoleAddHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}/role/{roleID}", userapiv1.UserRoleDeleteHandler).Methods("DELETE")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func main() {
	// Read command line args
	const defaultConfigFilePath = "./config.yaml"
	configFilePath := flag.String("config", defaultConfigFilePath, "file name of config.yaml")
	flag.Parse()

	// Read configure
	configAbsFilePath, _ := filepath.Abs(*configFilePath)
	cfg, err := config.InitConfig(configAbsFilePath)
	if err != nil {
		fmt.Printf("Failed to set config: %v", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %v", *cfg)

	// Setup API
	r := mux.NewRouter()
	setAPI(r)

	// Run Server
	corsObj := handlers.AllowedOrigins([]string{"*"})
	addr := fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.Port)
	logger.Info("start server with %s", addr)
	if err := http.ListenAndServe(addr, handlers.CORS(corsObj)(r)); err != nil {
		logger.Error("Failed to run server: %v", err)
		os.Exit(1)
	}
}
