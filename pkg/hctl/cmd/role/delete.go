package role

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var deleteRoleCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Role",
	Long:  "Delete role from the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		if err := handler.RoleDelete(projectName, name); err != nil {
			print.Fatal("Failed to delete the name %s from %s: %v", name, projectName, err)
		}

		print.Print("Role %s successfully deleted", name)
	},
}

func init() {
	deleteRoleCmd.Flags().String("project", "", "[Required] name of the project to which the role belongs")
	deleteRoleCmd.Flags().String("name", "", "[Required] name of the role")
	deleteRoleCmd.MarkFlagRequired("project")
	deleteRoleCmd.MarkFlagRequired("name")
}
