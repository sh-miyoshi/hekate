package client

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
		print.Error("client command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return clientCmd
}
