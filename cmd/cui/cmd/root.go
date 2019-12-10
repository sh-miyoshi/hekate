package cmd

import (
	"github.com/sh-miyoshi/jwt-server/cmd/cui/output"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/create"
	"github.com/spf13/cobra"
	"os"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
)

var serverAddr string
var outputFormat string
var isDebug bool

func init() {
	const defaultServerAddr = "http://localhost:8080"
	cobra.OnInitialize(initOutput)

	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", defaultServerAddr, "The address of jwt-server")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format: json, text")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "Output debug log")

	rootCmd.AddCommand(create.GetCreateCommand())
}

func initOutput() {
	if err := output.Init(outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to set output option: %v\n", err)
		os.Exit(1)
	}

	logger.InitLogger(isDebug, "")
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
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		os.Exit(1)
	}
}
