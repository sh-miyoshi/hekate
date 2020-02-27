package delete

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/spf13/cobra"
)

var deleteProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Delete Project",
	Long:  "Delete Project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		if err := handler.ProjectDelete(projectName); err != nil {
			fmt.Printf("Failed to delete project %s: %v\n", projectName, err)
			os.Exit(1)
		}

		fmt.Printf("Project %s successfully deleted\n", projectName)
	},
}

func init() {
	deleteProjectCmd.Flags().StringP("name", "n", "", "[Required] name of delete project")
	deleteProjectCmd.MarkFlagRequired("name")
}
