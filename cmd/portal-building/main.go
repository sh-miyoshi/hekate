package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/sh-miyoshi/hekate/pkg/logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./index.html")
	if err != nil {
		return
	}

	d := map[string]string{}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	tpl.Execute(w, d)
}

func main() {
	var port int
	var bindAddr string
	var logfile string
	flag.IntVar(&port, "port", 3000, "port number of server")
	flag.StringVar(&bindAddr, "addr", "0.0.0.0", "bind address of server")
	flag.StringVar(&logfile, "logfile", "", "file path for log, output to STDOUT if empty")
	flag.Parse()

	if err := logger.InitLogger(true, logfile); err != nil {
		fmt.Printf("Failed to init logger: %v", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	addr := fmt.Sprintf("%s:%d", bindAddr, port)
	logger.Info("Run server as %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Error("Failed to run server: %v", err)
		os.Exit(1)
	}
}
