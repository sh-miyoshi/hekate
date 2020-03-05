package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	clientCmd.AddCommand(addClientCmd)
	// clientCmd.AddCommand(deleteClientCmd)
	// clientCmd.AddCommand(getClientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Manage client in the project",
	Long:  `Manage client in the project`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("client command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return clientCmd
}
