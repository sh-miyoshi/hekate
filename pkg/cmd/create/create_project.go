package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/projectapi/v1"
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type projectInfo struct {
	Name                 string
	AccessTokenLifeSpan  int
	RefreshTokenLifeSpan int
}

var project projectInfo

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create New Project",
	Long:  "Create New Project",
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
		httpReq.Header.Add("Authorization", fmt.Sprintf("bearer %s", secret.AccessToken))
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

			// TODO(print)
			fmt.Printf("result: %v", res)
		}
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "[Required] set a name of new project")
	createProjectCmd.Flags().BoolVar(&project.Enabled, "enable", true, "set project enable")
	createProjectCmd.Flags().IntVar(&project.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	createProjectCmd.Flags().IntVar(&project.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	createProjectCmd.MarkFlagRequired("name")
}
