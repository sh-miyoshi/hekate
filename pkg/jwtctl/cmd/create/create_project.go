package create

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/spf13/cobra"
)

var project projectapi.ProjectCreateRequest

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create New Project",
	Long:  "Create New Project",
	Run: func(cmd *cobra.Command, args []string) {
		fileName, _ := cmd.Flags().GetString("file")

		if fileName == "" && project.Name == "" {
			fmt.Printf("\"name\" or \"file\" flag must be required")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		req := &project

		if fileName != "" {
			bytes, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Printf("Failed to read file %s: %v\n", fileName, err)
				os.Exit(1)
			}
			if err := json.Unmarshal(bytes, req); err != nil {
				fmt.Printf("Failed to parse json: %v\n", err)
				os.Exit(1)
			}
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		res, err := handler.ProjectAdd(req)
		if err != nil {
			fmt.Printf("Failed to add project %s: %v\n", project.Name, err)
			os.Exit(1)
		}

		format := output.NewProjectInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "name of new project")
	createProjectCmd.Flags().UintVar(&project.TokenConfig.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	createProjectCmd.Flags().UintVar(&project.TokenConfig.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	createProjectCmd.Flags().StringVar(&project.TokenConfig.SigningAlgorithm, "signAlg", "RS256", "token sigining algorithm, one of RS256, ")
	createProjectCmd.Flags().StringP("file", "f", "", "json file name of project info")
}
