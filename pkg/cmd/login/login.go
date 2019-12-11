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
		res, err := client.Do(httpReq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to request server: %v", err)
			return
		}
		defer res.Body.Close()

		// TODO(print res)
		fmt.Printf("result: %v\n", res)
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
