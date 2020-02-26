package create

import (
	"fmt"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/spf13/cobra"
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

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		req := &projectapi.ProjectCreateRequest{
			Name: project.Name,
			TokenConfig: projectapi.TokenConfig{
				AccessTokenLifeSpan:  project.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: project.RefreshTokenLifeSpan,
				SigningAlgorithm:     "RS256", // TODO(set param)
			},
		}

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
	createProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "[Required] name of new project")
	createProjectCmd.Flags().UintVar(&project.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	createProjectCmd.Flags().UintVar(&project.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	createProjectCmd.MarkFlagRequired("name")
}
