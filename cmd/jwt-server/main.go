package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/jwt-server/cmd/jwt-server/config"
	"github.com/sh-miyoshi/jwt-server/pkg/db"
	"github.com/sh-miyoshi/jwt-server/pkg/db/model"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	oidcapiv1 "github.com/sh-miyoshi/jwt-server/pkg/oidcapi/v1"
	projectapiv1 "github.com/sh-miyoshi/jwt-server/pkg/projectapi/v1"
	defaultrole "github.com/sh-miyoshi/jwt-server/pkg/role"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	userapiv1 "github.com/sh-miyoshi/jwt-server/pkg/userapi/v1"
	"github.com/sh-miyoshi/jwt-server/pkg/util"
	"net/http"
	"os"
	"time"
)

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// OpenID Connect API
	r.HandleFunc(basePath+"/project/{projectName}/.well-known/openid-configuration", oidcapiv1.ConfigGetHandler).Methods("GET")

	// Token API
	// TODO(depricated)
	r.HandleFunc(basePath+"/project/{projectName}/token", tokenapiv1.TokenCreateHandler).Methods("POST")

	// Project API
	r.HandleFunc(basePath+"/project", projectapiv1.AllProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project", projectapiv1.ProjectCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}", projectapiv1.ProjectDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}", projectapiv1.ProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}", projectapiv1.ProjectUpdateHandler).Methods("PUT")

	// User API
	r.HandleFunc(basePath+"/project/{projectName}/user", userapiv1.AllUserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/user", userapiv1.UserCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", userapiv1.UserDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", userapiv1.UserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", userapiv1.UserUpdateHandler).Methods("PUT")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/role/{roleID}", userapiv1.UserRoleAddHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/role/{roleID}", userapiv1.UserRoleDeleteHandler).Methods("DELETE")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func initDB(dbType, connStr, adminName, adminPassword string) error {
	if err := db.InitDBManager(dbType, connStr); err != nil {
		return errors.Wrap(err, "Failed to init database manager")
	}

	// Set Master Project if not exsits
	err := db.GetInst().ProjectAdd(&model.ProjectInfo{
		Name:      "master",
		CreatedAt: time.Now(),
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  5 * 60,            // 5 minutes, TODO(use const variable)
			RefreshTokenLifeSpan: 14 * 24 * 60 * 60, // 14 days, TODO(use const variable)
		},
	})
	if err != nil {
		if err == model.ErrProjectAlreadyExists {
			logger.Info("Master Project is already exists.")
		} else {
			return errors.Wrap(err, "Failed to create master project")
		}
	}

	err = db.GetInst().UserAdd(&model.UserInfo{
		ID:           uuid.New().String(),
		ProjectName:  "master",
		Name:         adminName,
		CreatedAt:    time.Now(),
		PasswordHash: util.CreateHash(adminPassword),
		Roles:        defaultrole.GetInst().GetList(), // set all roles
	})

	if err != nil {
		if err == model.ErrUserAlreadyExists {
			logger.Info("Admin user is already exists.")
		} else {
			return errors.Wrap(err, "Failed to create admin user")
		}
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

	// Initialize Default Role Handler
	if err := defaultrole.InitHandler(); err != nil {
		logger.Error("Failed to initialize default role handler: %+v", err)
		os.Exit(1)
	}

	// Initialize Token Config
	token.InitConfig(cfg.TokenIssuer, cfg.TokenSecretKey)

	// Initalize Database
	if err := initDB(cfg.DB.Type, cfg.DB.ConnectionString, cfg.AdminName, cfg.AdminPassword); err != nil {
		logger.Error("Failed to initialize database: %+v", err)
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
		logger.Error("Failed to run server: %+v", err)
		os.Exit(1)
	}
}
