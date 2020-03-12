package role

import (
	"github.com/spf13/cobra"
)

var addRoleCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New Role",
	Long:  "Add new role into the project",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}
