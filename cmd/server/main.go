package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/cmd/server/config"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	projectapiv1 "github.com/sh-miyoshi/jwt-server/pkg/projectapi/v1"
	roleapiv1 "github.com/sh-miyoshi/jwt-server/pkg/roleapi/v1"
	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	userapiv1 "github.com/sh-miyoshi/jwt-server/pkg/userapi/v1"
	"net/http"
	"os"
	"github.com/google/uuid"
)

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// Token API
	r.HandleFunc(basePath+"/token", tokenapiv1.TokenCreateHandler).Methods("POST")

	// Project API
	r.HandleFunc(basePath+"/project", projectapiv1.AllProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project", projectapiv1.ProjectCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}", projectapiv1.ProjectDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectID}", projectapiv1.ProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}", projectapiv1.ProjectUpdateHandler).Methods("PUT")

	// User API
	r.HandleFunc(basePath+"/project/{projectID}/user", userapiv1.AllUserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/user", userapiv1.UserCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}", userapiv1.UserUpdateHandler).Methods("PUT")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}/role/{roleID}", userapiv1.UserRoleAddHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}/user/{userID}/role/{roleID}", userapiv1.UserRoleDeleteHandler).Methods("DELETE")

	// Role API
	r.HandleFunc(basePath+"/project/{projectID}/role", roleapiv1.AllRoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/role", roleapiv1.RoleCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectID}/role/{roleID}", roleapiv1.RoleDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectID}/role/{roleID}", roleapiv1.RoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectID}/role/{roleID}", roleapiv1.RoleUpdateHandler).Methods("PUT")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func initDB(dbType, connStr, adminName, adminPassword string) error {
	if err := db.InitDBManager(dbType, connStr); err != nil {
		return err
	}

	// Set Master Project if not exsits
	err := db.GetInst().Project.Add(&model.ProjectInfo{
		ID:        "master",
		Name:      "master",
		Enabled:   true,
		CreatedAt: "now", // TODO
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  0, // TODO
			RefreshTokenLifeSpan: 0, // TODO
		},
	})
	if err != nil {
		logger.Info("Master Project is already exists. So nothing to do.")
		return nil
	}

	err = db.GetInst().User.Add(&model.UserInfo{
		ID:           uuid.New().String(),
		ProjectID:    "master",
		Name:         adminName,
		Enabled:      true,
		CreatedAt:    "now",         // TODO
		PasswordHash: util.CreateHash(adminPassword),
		Roles:        []string{},    // TODO
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Read command line args
	const defaultConfigFilePath = "./config.yaml"
	configFilePath := flag.String("config", defaultConfigFilePath, "file name of config.yaml")
	flag.Parse()

	// Read configure
	cfg, err := config.InitConfig(*configFilePath)
	if err != nil {
		fmt.Printf("Failed to set config: %v", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %v", *cfg)

	// Initalize Database
	if err := initDB("local", "./db", cfg.AdminName, cfg.AdminPassword); err != nil {
		logger.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}

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
