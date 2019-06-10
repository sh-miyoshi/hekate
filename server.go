package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type flagConfig struct {
	Port      int
	BindAddr  string
	LogFile   string
	ModeDebug bool
}

var config flagConfig

func parseCmdlineArgs() {
	const DefaultPort = 8080
	const DefaultBindAddr = "0.0.0.0"

	flag.IntVar(&config.Port, "port", DefaultPort, "set port number for server")
	flag.StringVar(&config.BindAddr, "bind", DefaultBindAddr, "set bind address for server")
	flag.StringVar(&config.LogFile, "logfile", "", "write log to file, output os.Stdout when do not set this")
	flag.BoolVar(&config.ModeDebug, "debug", false, "if true, run server as debug mode")
	flag.Parse()
}

func setAPI(r *mux.Router) {
	const basePath = "/api/v1"

	// TODO Add API

	// Health Check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}).Methods("GET")
}

func main() {
	parseCmdlineArgs()

	r := mux.NewRouter()
	setAPI(r)

	corsObj := handlers.AllowedOrigins([]string{"*"})

	addr := fmt.Sprintf("%s:%d", config.BindAddr, config.Port)
	if err := http.ListenAndServe(addr, handlers.CORS(corsObj)(r)); err != nil {
		os.Exit(1)
	}
}