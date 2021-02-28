package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/auth/v1/oidc"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	requestCmd.Flags().StringP("method", "X", "GET", "http method(GET, POST, DELETE, PUT)")
	requestCmd.Flags().StringP("input", "i", "", "input data")
	requestCmd.Flags().StringP("path", "p", "", "url path")

	requestCmd.MarkFlagRequired("path")
}

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "request to server",
	Long:  "request to server",
	Run: func(cmd *cobra.Command, args []string) {
		// get params
		serverAddr, _ := cmd.Flags().GetString("server")
		input, _ := cmd.Flags().GetString("input")
		path, _ := cmd.Flags().GetString("path")
		method, _ := cmd.Flags().GetString("method")

		if len(path) > 0 && path[0] != '/' {
			path = "/" + path
		}

		if len(input) > 0 && input[0] == '@' {
			print.Error("input file is not supported yet.") // TODO
			os.Exit(1)
		}

		// get token
		var rawJSON oidcapi.TokenResponse
		buf, err := ioutil.ReadFile("tmp/token.json")
		if err != nil {
			print.Error("Failed to read token.json: %v", err)
			os.Exit(1)
		}
		if err := json.Unmarshal(buf, &rawJSON); err != nil {
			print.Error("Failed to parse token.json: %v", err)
			os.Exit(1)
		}
		tkn := rawJSON.AccessToken

		// request to server
		u := serverAddr + path
		req, err := http.NewRequest(method, u, nil)
		if err != nil {
			print.Error("Failed to request to server: %v", err)
			os.Exit(1)
		}
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", tkn))
		if input != "" {
			req.Header.Add("Content-Type", "application/json")
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					ServerName:         req.Host,
					InsecureSkipVerify: true,
				},
			},
		}

		dump, _ := httputil.DumpRequest(req, false)
		print.Debug("server request dump: %q", dump)

		// get, parse and show response
		res, err := client.Do(req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				print.Print("No response")
				return
			}

			print.Error("Failed to get response: %v", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		dump, _ = httputil.DumpResponse(res, false)
		print.Debug("server response dump: %q", dump)

		io.Copy(os.Stdout, res.Body)
	},
}
