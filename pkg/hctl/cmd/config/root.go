package config

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(setConfigCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration of hctl command",
	Long:  `Manage configuration of hctl command`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		print.Error("config command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return configCmd
}
