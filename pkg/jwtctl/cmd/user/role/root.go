package role

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	roleCmd.AddCommand(addRoleCmd)
}

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage role in the user",
	Long:  `Manage role in the user`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("role command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return roleCmd
}
