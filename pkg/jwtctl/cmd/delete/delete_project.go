package delete

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var projectName string

var deleteProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Delete Project",
	Long:  "Delete Project",
	Run: func(cmd *cobra.Command, args []string) {
		secret, err := config.GetSecretToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

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
