package main

import (
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/cmd/hekate/config"
	auditapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/audit"
	clientapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/client"
	roleapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
	oidcapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
	projectapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	sessionapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/session"
	userapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/user"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/oidc"
	"github.com/sh-miyoshi/hekate/pkg/oidc/token"
	defaultrole "github.com/sh-miyoshi/hekate/pkg/role"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s: %s called", r.Method, r.URL.String())

		vars := mux.Vars(r)
		projectName := vars["projectName"]

		// Validate project name
		if projectName != "" {
			if _, err := db.GetInst().ProjectGet(projectName); err != nil {
				if errors.Contains(err, model.ErrNoSuchProject) || errors.Contains(err, model.ErrProjectValidateFailed) {
					errors.PrintAsInfo(errors.Append(err, "No such project %s", projectName))
					http.Error(w, "Project Not Found", http.StatusNotFound)
				} else {
					errors.Print(errors.Append(err, "Failed to validate project name"))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
func setAPI(r *mux.Router, cfg *config.GlobalConfig) {
	const basePath = "/api/v1"

	// OpenID Connect API
	r.HandleFunc(basePath+"/project/{projectName}/.well-known/openid-configuration", oidcapiv1.ConfigGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/token", oidcapiv1.TokenHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/certs", oidcapiv1.CertsHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/auth", oidcapiv1.AuthGETHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/auth", oidcapiv1.AuthPOSTHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/userinfo", oidcapiv1.UserInfoHandler).Methods("GET", "POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/revoke", oidcapiv1.RevokeHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/login", oidcapiv1.UserLoginHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/consent", oidcapiv1.ConsentHandler).Methods("POST")

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
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/change-password", userapiv1.UserChangePasswordHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/logout", userapiv1.UserLogoutHandler).Methods("POST")

	// Client API
	r.HandleFunc(basePath+"/project/{projectName}/client", clientapiv1.AllClientGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/client", clientapiv1.ClientCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", clientapiv1.ClientDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", clientapiv1.ClientGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", clientapiv1.ClientUpdateHandler).Methods("PUT")

	// Custom Role API
	r.HandleFunc(basePath+"/project/{projectName}/role", roleapiv1.AllRoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/role", roleapiv1.RoleCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", roleapiv1.RoleDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", roleapiv1.RoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", roleapiv1.RoleUpdateHandler).Methods("PUT")

	// Session API
	r.HandleFunc(basePath+"/project/{projectName}/session/{sessionID}", sessionapiv1.SessionDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/session/{sessionID}", sessionapiv1.SessionGetHandler).Methods("GET")

	// Audit API
	r.HandleFunc(basePath+"/project/{projectName}/audit", auditapiv1.AuditGetHandler).Methods("GET")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.GetInst().Ping(); err != nil {
			http.Error(w, "DB Ping Failed", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	}).Methods("GET")

	// File Server for User Login Page
	fsCSS := http.FileServer(http.Dir(path.Join(cfg.UserLoginResourceDir, "/css")))
	pt := path.Join(authCodeUserLoginResourcePath, "/css") + "/"
	r.PathPrefix(pt).Handler(http.StripPrefix(pt, fsCSS))
	fsImg := http.FileServer(http.Dir(path.Join(cfg.UserLoginResourceDir, "/img")))
	pt = path.Join(authCodeUserLoginResourcePath, "/img") + "/"
	r.PathPrefix(pt).Handler(http.StripPrefix(pt, fsImg))

	r.Use(loggingMiddleware)
}

func initDB(dbType, connStr, adminName, adminPassword string) *errors.Error {
	if err := db.InitDBManager(dbType, connStr); err != nil {
		return errors.Append(err, "Failed to init database manager")
	}

	// Set Master Project if not exsits
	err := db.GetInst().ProjectAdd(&model.ProjectInfo{
		Name:         "master",
		CreatedAt:    time.Now(),
		PermitDelete: false,
		TokenConfig: &model.TokenConfig{
			AccessTokenLifeSpan:  model.DefaultAccessTokenExpiresTimeSec,
			RefreshTokenLifeSpan: model.DefaultRefreshTokenExpiresTimeSec,
			SigningAlgorithm:     "RS256",
		},
		AllowGrantTypes: []model.GrantType{
			model.GrantTypeAuthorizationCode,
			model.GrantTypeClientCredentials,
			model.GrantTypeRefreshToken,
			model.GrantTypePassword, // TODO(for debug)
		},
	})
	if err != nil {
		if errors.Contains(err, model.ErrProjectAlreadyExists) {
			logger.Info("Master Project is already exists.")
		} else {
			return errors.Append(err, "Failed to create master project")
		}
	} else {
		logger.Debug("Add master project")
	}

	err = db.GetInst().UserAdd("master", &model.UserInfo{
		ID:           uuid.New().String(),
		ProjectName:  "master",
		Name:         adminName,
		CreatedAt:    time.Now(),
		PasswordHash: util.CreateHash(adminPassword),
		SystemRoles: []string{
			// append cluster admin role
			"read-cluster",
			"write-cluster",
		},
	})

	if err != nil {
		if errors.Contains(err, model.ErrUserAlreadyExists) {
			logger.Info("Admin user is already exists.")
		} else {
			return errors.Append(err, "Failed to create admin user")
		}
	} else {
		logger.Debug("Add admin user to master project")
	}

	return nil
}

func initAll(cfg *config.GlobalConfig) *errors.Error {
	// Initialize logger
	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %v", *cfg)

	// Check login resource directory struct
	if err := cfg.CheckLoginResDirStruct(); err != nil {
		return errors.Append(err, "Login resource directory is broken")
	}

	// Initialize Default Role Handler
	if err := defaultrole.InitHandler(); err != nil {
		return errors.Append(err, "Failed to initialize default role handler")
	}

	// Initialize Token Config
	token.InitConfig(cfg.HTTPSConfig.Enabled)

	// Initialize OIDC Config
	oidc.InitConfig(cfg.HTTPSConfig.Enabled, cfg.AuthCodeExpiresTime, cfg.UserLoginResourceDir, authCodeUserLoginResourcePath)

	// Initalize Database
	if err := initDB(cfg.DB.Type, cfg.DB.ConnectionString, cfg.AdminName, cfg.AdminPassword); err != nil {
		return errors.Append(err, "Failed to initialize database")
	}

	// Initialize Audit Events Database
	typ := cfg.AuditDB.Type
	connStr := cfg.AuditDB.ConnectionString
	if typ == "" {
		typ = cfg.DB.Type
		connStr = cfg.DB.ConnectionString
	}
	if err := audit.Init(typ, connStr); err != nil {
		return errors.Append(err, "Failed to initialize audit events database")
	}

	return nil
}
