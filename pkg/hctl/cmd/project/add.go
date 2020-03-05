package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/spf13/cobra"
)

var project projectapi.ProjectCreateRequest

var addProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New Project",
	Long:  "Add New Project",
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
				fmt.Printf("Failed to parse input file to json: %v\n", err)
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
	addProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "name of new project")
	addProjectCmd.Flags().UintVar(&project.TokenConfig.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	addProjectCmd.Flags().UintVar(&project.TokenConfig.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	addProjectCmd.Flags().StringVar(&project.TokenConfig.SigningAlgorithm, "signAlg", "RS256", "token sigining algorithm, one of RS256, ")
	addProjectCmd.Flags().StringP("file", "f", "", "json file name of project info")
}
