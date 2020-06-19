package project

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var deleteProjectCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Project",
	Long:  "Delete Project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		if err := handler.ProjectDelete(projectName); err != nil {
			print.Fatal("Failed to delete project %s: %v", projectName, err)
		}

		print.Print("Project %s successfully deleted", projectName)
	},
}

func init() {
	deleteProjectCmd.Flags().StringP("name", "n", "", "[Required] name of delete project")
	deleteProjectCmd.MarkFlagRequired("name")
}
