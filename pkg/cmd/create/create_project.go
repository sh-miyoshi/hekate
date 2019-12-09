package create

import (
	"fmt"
	"github.com/spf13/cobra"
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
