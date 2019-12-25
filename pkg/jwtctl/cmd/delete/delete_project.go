package delete

import (
	"encoding/json"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var projectName string

var deleteProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Delete Project",
	Long:  "Delete Project",
	Run: func(cmd *cobra.Command, args []string) {
		// Get Secret Info
		secretFile := filepath.Join(config.Get().ConfigDir, "secret")
		buf, err := ioutil.ReadFile(secretFile)
		if err != nil {
			fmt.Printf("Failed to read secret file: %v\n", err)
			fmt.Println("You need to `jwtctl login` at first")
			os.Exit(1)
		}

		var secret tokenapi.TokenResponse
		if err = json.Unmarshal(buf, &secret); err != nil {
			fmt.Printf("Failed to parse secret json: %v", err)
			os.Exit(1)
		}

		// TODO(Validate secret)

		serverAddr := config.Get().ServerAddr
		url := fmt.Sprintf("%s/api/v1/project/%s", serverAddr, projectName)
		httpReq, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			logger.Error("<Program Bug> Failed to create http request: %v", err)
			os.Exit(1)
		}
		httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", secret.AccessToken))
		client := &http.Client{}
		httpRes, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("Failed to request server: %v", err)
			os.Exit(1)
		}
		defer httpRes.Body.Close()

		switch httpRes.StatusCode {
		case 204:
			fmt.Printf("Successfully deleted\n")
		}
	},
}

func init() {
	deleteProjectCmd.Flags().StringVarP(&projectName, "name", "n", "", "[Required] set a name of delete project")
	deleteProjectCmd.MarkFlagRequired("name")
}
