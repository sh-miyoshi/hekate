package cmd

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tokenCmd)
}

var rootCmd = &cobra.Command{
	Use:   "debug-api",
	Short: "debug-api is a debuging tool for hekate server api",
	Long:  "debug-api is a debuging tool for hekate server api",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute method run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		print.Error("Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}
