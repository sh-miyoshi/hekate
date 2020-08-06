package main

import (
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/cmd/hekate/config"
	clientapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/client"
	roleapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
	oidcapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
	projectapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	sessionapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/session"
	userapiv1 "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/user"
	"github.com/sh-miyoshi/hekate/pkg/db"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/errors"
	"github.com/sh-miyoshi/hekate/pkg/logger"
	"github.com/sh-miyoshi/hekate/pkg/util"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("%s: %s called", r.Method, r.URL.String())
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
