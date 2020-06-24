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
	addRoleCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	addRoleCmd.Flags().String("user", "", "[Required] name of user")
	addRoleCmd.Flags().StringP("role", "r", "", "[Required] role name to add to the user")
	addRoleCmd.Flags().String("type", "system", "role type (system or custom)")

	addRoleCmd.MarkFlagRequired("project")
	addRoleCmd.MarkFlagRequired("user")
	addRoleCmd.MarkFlagRequired("role")
}

var addRoleCmd = &cobra.Command{
	Use:   "add",
	Short: "Add role to the user",
	Long:  "Add role to the user",
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
		if err := handler.UserRoleAdd(projectName, userName, roleName, roleType); err != nil {
			print.Fatal("Failed to add role %s to user %s in %s: %v", roleName, userName, projectName, err)
		}

		print.Print("Successfully added")
	},
}
