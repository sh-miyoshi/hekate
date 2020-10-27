package secret

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var projectSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage project secret",
	Long:  "Manage project secret",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		print.Print("secret command requires subcommand")
	},
}

func init() {
	projectSecretCmd.AddCommand(getCmd)
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return projectSecretCmd
}
