package role

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	deleteRoleCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	deleteRoleCmd.Flags().String("user", "", "name of user")
	deleteRoleCmd.Flags().StringP("role", "r", "", "role name to add to the user")
	deleteRoleCmd.Flags().String("type", "system", "role type (system or custom)")

	deleteRoleCmd.MarkFlagRequired("project")
	deleteRoleCmd.MarkFlagRequired("user")
	deleteRoleCmd.MarkFlagRequired("role")
}

var deleteRoleCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete rol from the user",
	Long:  "Delete rol from the user",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("user")
		roleName, _ := cmd.Flags().GetString("role")
		typ, _ := cmd.Flags().GetString("type")

		roleType := model.RoleSystem
		switch typ {
		case "system":
			roleType = model.RoleSystem
		case "custom":
			roleType = model.RoleCustom
		default:
			print.Error("Please set role type to system or custom.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		if err := handler.UserRoleDelete(projectName, userName, roleName, roleType); err != nil {
			print.Fatal("Failed to delete role %s from user %s in %s: %v", roleName, userName, projectName, err)
		}

		print.Print("Successfully deleted")
	},
}
