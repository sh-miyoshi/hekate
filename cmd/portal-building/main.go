package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
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
	flag.IntVar(&port, "port", 3000, "port number of server")
	flag.StringVar(&bindAddr, "addr", "0.0.0.0", "bind address of server")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")

	addr := fmt.Sprintf("%s:%d", bindAddr, port)
	if err := http.ListenAndServe(addr, r); err != nil {
		os.Exit(1)
	}
}
