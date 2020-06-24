package role

import (
	"github.com/spf13/cobra"
)

var getRoleCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Roles in the cluster",
	Long:  "Get Roles in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(implement this)
	},
}

func init() {
}
