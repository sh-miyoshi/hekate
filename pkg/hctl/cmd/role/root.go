package role

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	roleCmd.AddCommand(addRoleCmd)
	roleCmd.AddCommand(deleteRoleCmd)
	roleCmd.AddCommand(getRoleCmd)
}

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage custom role in the project",
	Long:  `Manage custom role in the project`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		print.Error("role command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return roleCmd
}
