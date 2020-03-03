package role

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addRoleCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	addRoleCmd.Flags().String("user", "", "name of user")
	addRoleCmd.Flags().StringSliceP("roles", "r", nil, "role list to add to the user")
	addRoleCmd.Flags().String("type", "system", "role type (system or custom)")

	addRoleCmd.MarkFlagRequired("project")
	addRoleCmd.MarkFlagRequired("user")
}

var addRoleCmd = &cobra.Command{
	Use:   "add",
	Short: "Add role to the user",
	Long:  "Add role to the user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Not Implemented yet\n")
	},
}
