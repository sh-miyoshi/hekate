package cmd

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var debugMode bool

func init() {
	cobra.OnInitialize(initOutput)

	rootCmd.PersistentFlags().String("server", "http://localhost:18443", "server address")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Output debug message")

	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(requestCmd)
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

func initOutput() {
	print.Init(debugMode)
}
