package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type projectInfo struct {
	Name                 string
	AccessTokenLifeSpan  uint
	RefreshTokenLifeSpan uint
}

var project projectInfo

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create New Project",
	Long:  "Create New Project",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		serverAddr := config.Get().ServerAddr
		req := projectapi.ProjectCreateRequest{
			Name: project.Name,
			TokenConfig: &projectapi.TokenConfig{
				AccessTokenLifeSpan:  project.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: project.RefreshTokenLifeSpan,
			},
		}

		url := fmt.Sprintf("%s/api/v1/project", serverAddr)
		body, err := json.Marshal(req)
		if err != nil {
			logger.Error("<Program Bug> Failed to marshal JSON: %v", err)
			os.Exit(1)
		}

		httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			logger.Error("<Program Bug> Failed to create http request: %v", err)
			os.Exit(1)
		}
		httpReq.Header.Add("Content-Type", "application/json")
		httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
		client := &http.Client{}
		httpRes, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("Failed to request server: %v", err)
			os.Exit(1)
		}
		defer httpRes.Body.Close()

		switch httpRes.StatusCode {
		case 200:
			var res projectapi.ProjectGetResponse
			if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
				logger.Error("<Program Bug> Failed to parse http response: %v", err)
				os.Exit(1)
			}

			format := output.NewProjectInfoFormat(&res)
			output.Print(format)
		}
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "[Required] set a name of new project")
	createProjectCmd.Flags().UintVar(&project.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	createProjectCmd.Flags().UintVar(&project.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	createProjectCmd.MarkFlagRequired("name")
}
