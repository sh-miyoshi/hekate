package apiclient

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

// Handler ...
type Handler struct {
	client      *http.Client
	serverAddr  string
	accessToken string
}

// NewHandler ...
func NewHandler(serverAddr string, accessToken string) *Handler {
	h := &Handler{
		serverAddr:  serverAddr,
		accessToken: accessToken,
	}

	// TODO(set correct params)
	insecure := true
	timeout := time.Duration(10 * time.Second)

	h.client = createClient(serverAddr, insecure, timeout)

	return h
}

// trimHTTPPrefix trims "http://" and "https://"
func trimHTTPPrefix(addr string) string {
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	return addr
}

func createClient(serverAddr string, insecure bool, timeout time.Duration) *http.Client {
	tlsConfig := tls.Config{
		ServerName: trimHTTPPrefix(serverAddr),
	}

	if insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return client
}
