package cmd

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/create"
	"github.com/spf13/cobra"
	"os"
)

var serverAddr string

func init() {
	const defaultServerAddr = "http://localhost:8080"

	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", defaultServerAddr, "The address of jwt-server")

	rootCmd.AddCommand(create.GetCreateCommand())
}

var rootCmd = &cobra.Command{
	Use:   "jwt",
	Short: "jwt is a command line tool for jwt-server",
	Long:  "jwt is a command line tool for jwt-server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute method run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
