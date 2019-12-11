package login

import (
	"github.com/spf13/cobra"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"path/filepath"
	"io/ioutil"
	"os"
)

var (
	userName    string
	password    string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr := config.Get().ServerAddr
		project := config.Get().ProjectName
		logger.Debug("server address: %s", serverAddr)

		// TODO(input password in STDIN)

		req := tokenapi.TokenRequest{
			Name: userName,
			Secret: password,
			AuthType: "password",
		}

		url := fmt.Sprintf("%s/api/v1/project/%s/token", serverAddr, project)
		body, err := json.Marshal(req)
		if err != nil {
			logger.Error("<Program Bug> Failed to marshal JSON: %v", err)
			return
		}

		httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			logger.Error("<Program Bug> Failed to create http request: %v", err)
			return
		}
		httpReq.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		httpRes, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("Failed to request server: %v", err)
			return
		}
		defer httpRes.Body.Close()

		switch httpRes.StatusCode {
		case 200:
			var res tokenapi.TokenResponse
			if err := json.NewDecoder(httpRes.Body).Decode(&res);err!= nil{
				logger.Error("<Program Bug> Failed to parse http response: %v", err)
				return
			}
			// Output to secret file
			secretFile := filepath.Join(config.Get().ConfigDir,"secret")
			bytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				logger.Error("<Program Bug> Failed to marshal json: %v", err)
				return
			}

			ioutil.WriteFile(secretFile, bytes, os.ModePerm)
			fmt.Println("Successfully logged in")
		case 401, 404:
			fmt.Println("Failed to login to system")
			fmt.Println("Please cheak user name or password (or project name)")
		case 500:
			fmt.Println("Internal Server Error is occured")
			fmt.Println("Please contact to your server administrator")
		default:
			logger.Error("<Program Bug> Unexpected http response code: %d", httpRes.StatusCode)
		}
	},
}

func init() {
	loginCmd.Flags().StringVarP(&userName, "name", "n", "", "Login User Name")
	loginCmd.Flags().StringVar(&password, "password", "", "Login User Password")
	loginCmd.MarkFlagRequired("name")
}

// GetLoginCommand ...
func GetLoginCommand() *cobra.Command {
	return loginCmd
}
