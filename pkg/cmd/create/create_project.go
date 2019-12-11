package create

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	projectapi "github.com/sh-miyoshi/jwt-server/pkg/projectapi/v1"
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type projectInfo struct {
	Name                 string
	Enabled              bool
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
			fmt.Println("You need to `jwt login` at first")
			os.Exit(1)
		}

		var secret tokenapi.TokenResponse
		if err = yaml.Unmarshal(buf, &secret); err != nil {
			logger.Error("<Program Bug> Failed to marshal secret yaml: %v", err)
			os.Exit(1)
		}

		serverAddr := config.Get().ServerAddr
		logger.Debug("server address: %s", serverAddr)

		req := projectapi.ProjectCreateRequest{
			Name:    project.Name,
			Enabled: project.Enabled,
			TokenConfig: &projectapi.TokenConfig{
				AccessTokenLifeSpan:  project.AccessTokenLifeSpan,
				RefreshTokenLifeSpan: project.RefreshTokenLifeSpan,
			},
		}
		logger.Debug("project create request: %v", req)

		fmt.Println("create new project")
	},
}

func init() {
	createProjectCmd.Flags().StringVarP(&project.Name, "name", "n", "", "[Required] set a name of new project")
	createProjectCmd.Flags().BoolVar(&project.Enabled, "enable", true, "set project enable")
	createProjectCmd.Flags().IntVar(&project.AccessTokenLifeSpan, "accessExpires", 5*60, "access token life span [sec]")
	createProjectCmd.Flags().IntVar(&project.RefreshTokenLifeSpan, "refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	createProjectCmd.MarkFlagRequired("name")
}
