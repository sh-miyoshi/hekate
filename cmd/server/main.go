package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/sh-miyoshi/jwt-server/pkg/token"
	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
)

type globalConfig struct {
	Port                int
	BindAddr            string
	LogFile             string
	ModeDebug           bool
	AdminName           string
	AdminPassword       string
	TokenExpiredTimeSec int
	TokenIssuer         string
}

var config globalConfig

func parseCmdlineArgs() {
	const DefaultPort = 8080
	const DefaultBindAddr = "0.0.0.0"
	const DefaultAdminUser = "admin"
	const DefaultAdminPassword = "password"
	const DefaultTokenExpiredTime = 3600 // 1h
	const DefaultTokenIssuer = "jwt-server"

	flag.IntVar(&config.Port, "port", DefaultPort, "set port number for server")
	flag.StringVar(&config.BindAddr, "bind", DefaultBindAddr, "set bind address for server")
	flag.StringVar(&config.LogFile, "logfile", "", "write log to file, output os.Stdout when do not set this option")
	flag.BoolVar(&config.ModeDebug, "debug", false, "if true, run server as debug mode")
	flag.StringVar(&config.AdminName, "admin-name", DefaultAdminUser, "user name of system admin")
	flag.StringVar(&config.AdminPassword, "admin-password", DefaultAdminPassword, "password of system admin")
	flag.IntVar(&config.TokenExpiredTimeSec, "expired-time", DefaultTokenExpiredTime, "JWT token expired time [second]")
	flag.StringVar(&config.TokenIssuer, "issuer", DefaultTokenIssuer, "issuer of JWT token")
	flag.Parse()
}

func setEnvVar(key string, target *string) {
	val := os.Getenv(key)
	if len(val) > 0 {
		target = &val
	}
}

func parseOSEnvironment() {
	setEnvVar("JWT_SERVER_ADMIN_NAME", &config.AdminName)
	setEnvVar("JWT_SERVER_ADMIN_PASSWORD", &config.AdminPassword)
	setEnvVar("JWT_SERVER_TOKEN_ISSUER", &config.TokenIssuer)

	// set token rexired time
	if len(os.Getenv("JWT_SERVER_TOKEN_EXPIRED_TIME")) > 0 {
		time, _ := strconv.Atoi(os.Getenv("JWT_SERVER_TOKEN_EXPIRED_TIME"))
		if time > 0 {
			config.TokenExpiredTimeSec = time
		}
	}
}

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// Add API
	r.HandleFunc(basePath+"/token", tokenapiv1.CreateTokenHandler).Methods("POST")

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func main() {
	parseOSEnvironment()
	parseCmdlineArgs()

	logger.InitLogger(config.ModeDebug, config.LogFile)

	// Initialize Token Config
	secretKey := uuid.New().String()
	expiredTime := time.Second * time.Duration(config.TokenExpiredTimeSec)
	token.InitConfig(expiredTime, config.TokenIssuer, secretKey)

	r := mux.NewRouter()
	setAPI(r)

	corsObj := handlers.AllowedOrigins([]string{"*"})

	addr := fmt.Sprintf("%s:%d", config.BindAddr, config.Port)
	logger.Info("start server with %s", addr)
	if err := http.ListenAndServe(addr, handlers.CORS(corsObj)(r)); err != nil {
		os.Exit(1)
	}
}
