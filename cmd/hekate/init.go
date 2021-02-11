package main

import (
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	adminauditapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/audit"
	adminauthnapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/authn"
	adminclientapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/client"
	adminroleapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/customrole"
	adminkeysapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/keys"
	adminoauthapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/oauth"
	adminoidcapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/oidc"
	adminprojectapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/project"
	adminsessionapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/session"
	adminuserapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/user"
	userapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/user/v1"
	"github.com/sh-miyoshi/hekate/pkg/audit"
	"github.com/sh-miyoshi/hekate/pkg/config"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/login"
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
					errors.WriteHTTPError(w, "Not Found", err, http.StatusNotFound)
				} else {
					errors.Print(errors.Append(err, "Failed to validate project name"))
					errors.WriteHTTPError(w, "Internal Server Error", err, http.StatusInternalServerError)
				}
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
func setAPI(r *mux.Router) {
	cfg := config.Get()

	//------------------------------
	// Admin APIs
	//------------------------------
	const basePath = "/adminapi/v1"

	// OpenID Connect API
	r.HandleFunc(basePath+"/project/{projectName}/.well-known/openid-configuration", adminoidcapiv1.ConfigGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/token", adminoidcapiv1.TokenHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/certs", adminoidcapiv1.CertsHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/auth", adminoidcapiv1.AuthGETHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/auth", adminoidcapiv1.AuthPOSTHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/userinfo", adminoidcapiv1.UserInfoHandler).Methods("GET", "POST")
	r.HandleFunc(basePath+"/project/{projectName}/openid-connect/revoke", adminoidcapiv1.RevokeHandler).Methods("POST")

	// OAuth
	r.HandleFunc(basePath+"/project/{projectName}/oauth/device", adminoauthapiv1.DeviceRegisterHandler).Methods("POST")

	// Authenticate API
	r.HandleFunc(basePath+"/project/{projectName}/authn/login", adminauthnapiv1.UserLoginHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/authn/consent", adminauthnapiv1.ConsentHandler).Methods("POST")

	// Project API
	r.HandleFunc(basePath+"/project", adminprojectapiv1.AllProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project", adminprojectapiv1.ProjectCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}", adminprojectapiv1.ProjectDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}", adminprojectapiv1.ProjectGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}", adminprojectapiv1.ProjectUpdateHandler).Methods("PUT")

	// Keys API
	r.HandleFunc(basePath+"/project/{projectName}/keys", adminkeysapiv1.KeysGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/keys/reset", adminkeysapiv1.KeysResetHandler).Methods("POST")

	// User API
	r.HandleFunc(basePath+"/project/{projectName}/user", adminuserapiv1.AllUserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/user", adminuserapiv1.UserCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", adminuserapiv1.UserDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", adminuserapiv1.UserGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}", adminuserapiv1.UserUpdateHandler).Methods("PUT")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/role/{roleID}", adminuserapiv1.UserRoleAddHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/role/{roleID}", adminuserapiv1.UserRoleDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/reset-password", adminuserapiv1.UserResetPasswordHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/user/{userID}/unlock", adminuserapiv1.UserUnlockHandler).Methods("POST")

	// Client API
	r.HandleFunc(basePath+"/project/{projectName}/client", adminclientapiv1.AllClientGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/client", adminclientapiv1.ClientCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", adminclientapiv1.ClientDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", adminclientapiv1.ClientGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/client/{clientID}", adminclientapiv1.ClientUpdateHandler).Methods("PUT")

	// Custom Role API
	r.HandleFunc(basePath+"/project/{projectName}/role", adminroleapiv1.AllRoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/role", adminroleapiv1.RoleCreateHandler).Methods("POST")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", adminroleapiv1.RoleDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", adminroleapiv1.RoleGetHandler).Methods("GET")
	r.HandleFunc(basePath+"/project/{projectName}/role/{roleID}", adminroleapiv1.RoleUpdateHandler).Methods("PUT")

	// Session API
	r.HandleFunc(basePath+"/project/{projectName}/session/{sessionID}", adminsessionapiv1.SessionDeleteHandler).Methods("DELETE")
	r.HandleFunc(basePath+"/project/{projectName}/session/{sessionID}", adminsessionapiv1.SessionGetHandler).Methods("GET")

	// Audit API
	r.HandleFunc(basePath+"/project/{projectName}/audit", adminauditapiv1.AuditGetHandler).Methods("GET")

	//------------------------------
	// User APIs
	//------------------------------
	r.HandleFunc("/userapi/v1/project/{projectName}/user/{userID}/change-password", userapiv1.ChangePasswordHandler)
	r.HandleFunc("/userapi/v1/project/{projectName}/user/{userID}/logout", userapiv1.LogoutHandler).Methods("POST")

	//------------------------------
	// Other APIs
	//------------------------------

	// Device Login HTML Page
	r.HandleFunc("/resource/project/{projectName}/devicelogin", adminoauthapiv1.DeviceLoginPageHandler).Methods("GET")
	r.HandleFunc("/resource/project/{projectName}/deviceverify", adminoauthapiv1.DeviceUserCodeVerifyHandler).Methods("POST")
	r.HandleFunc("/resource/project/{projectName}/devicecomplete", func(w http.ResponseWriter, r *http.Request) {
		login.WriteDeviceLoginCompletePage(w)
	}).Methods("GET")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.GetInst().Ping(); err != nil {
			errors.WriteHTTPError(w, "Ping Failed", err, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("ok"))
	}).Methods("GET")

	// File Server for User Login Page
	fs := http.FileServer(http.Dir(path.Join(cfg.UserLoginResourceDir, "/static")))
	pt := path.Join(cfg.LoginStaticResourceURL, "/static") + "/"
	r.PathPrefix(pt).Handler(http.StripPrefix(pt, fs))

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
			AccessTokenLifeSpan:  model.DefaultAccessTokenExpiresInSec,
			RefreshTokenLifeSpan: model.DefaultRefreshTokenExpiresInSec,
			SigningAlgorithm:     "RS256",
		},
		AllowGrantTypes: []model.GrantType{
			model.GrantTypeAuthorizationCode,
			model.GrantTypeClientCredentials,
			model.GrantTypeRefreshToken,
			model.GrantTypeDevice,
			model.GrantTypePassword, // TODO(for debug)
		},
		UserLock: model.UserLock{
			Enabled:          false,
			MaxLoginFailure:  model.DefaultMaxLoginFailure,
			LockDuration:     model.DefaultLockDuration,
			FailureResetTime: model.DefaultFailureResetTime,
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

func initAll() *errors.Error {
	cfg := config.Get()

	// Initialize logger
	logger.InitLogger(cfg.ModeDebug, cfg.LogFile)
	logger.Debug("Start with config: %+v", *cfg)

	// Initialize settings
	if err := defaultrole.InitHandler(); err != nil {
		return errors.Append(err, "Failed to initialize default role handler")
	}
	logger.Debug("Successfully initialize system role")

	// Initalize Database
	if err := initDB(cfg.DB.Type, cfg.DB.ConnectionString, cfg.AdminName, cfg.AdminPassword); err != nil {
		return errors.Append(err, "Failed to initialize database")
	}
	logger.Debug("Successfully initialize database")

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
	logger.Debug("Successfully initialize audit db with type: %s", typ)

	// Initialize DBGC
	db.InitGC(cfg.DBGCInterval)
	logger.Debug("Start database GC per %d [sec]", cfg.DBGCInterval)

	return nil
}
