package role

import (
	"github.com/spf13/cobra"
)

var deleteRoleCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Role",
	Long:  "Delete role from the project",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(implement this)
	},
}

func init() {
}
