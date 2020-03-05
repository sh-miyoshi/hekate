package user

import (
	"fmt"

	"github.com/sh-miyoshi/hekate/pkg/jwtctl/cmd/user/role"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(deleteUserCmd)
	userCmd.AddCommand(getUserCmd)
	userCmd.AddCommand(role.GetCommand())
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user in the project",
	Long:  `Manage user in the project`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("user command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return userCmd
}
