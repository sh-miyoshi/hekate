package role

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/db/model"
	"github.com/sh-miyoshi/hekate/pkg/jwtctl/config"

	"github.com/spf13/cobra"
)

func init() {
	addRoleCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	addRoleCmd.Flags().String("user", "", "name of user")
	addRoleCmd.Flags().StringP("role", "r", "", "role name to add to the user")
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
			fmt.Printf("Please set role type to system or custom.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		if err := handler.UserRoleAdd(projectName, userName, roleName, roleType); err != nil {
			fmt.Printf("Failed to add role %s to user %s in %s: %v", roleName, userName, projectName, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added\n")
	},
}
