package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var logger = log.New(os.Stderr, "[TESTSERVER]", log.LUTC|log.LstdFlags)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("%s method is approved", r.Method)

	// return success message
	w.Write([]byte(r.Method + " is success"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Hello service was called with")
	logger.Printf("Host: %s", r.Host)
	fmt.Fprint(w, "Hello")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Echo service was called : Content-Length [%d] bytes", r.ContentLength)

	length := fmt.Sprintf("%d", r.ContentLength)
	w.Header().Set("Content-Length", length)
	_, err := io.Copy(w, r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func main() {
	var logFile string
	var httpBindPort int // defualt port number is 10000
	flag.StringVar(&logFile, "logfile", "", "-logfile=<log-file-name>")
	flag.IntVar(&httpBindPort, "port", 10000, "-port=PortNumber")
	flag.Parse()
	HTTPBindAddr := "0.0.0.0:" + strconv.Itoa(httpBindPort)

	// output log to file
	if logFile != "" {
		file, err := os.Create(logFile)
		if err != nil {
			logger.Printf("failed to create logFile %v", logFile)
			os.Exit(1)
		}
		//do not call file.Close() because logger write log through file.Writer

		logger.SetOutput(file)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/hello", helloHandler).Methods("GET")
	r.HandleFunc("/echo", echoHandler).Methods("POST")

	logger.Printf("Start Server\n")
	if err := http.ListenAndServe(HTTPBindAddr, r); err != nil {
		logger.Printf("%v", err)
		os.Exit(1)
	}
}
