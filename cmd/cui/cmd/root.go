package cmd

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/cmd/cui/output"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/config"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/create"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	"github.com/spf13/cobra"

	"os"
)

var (
	outputFormat string
	configDir    string
)

func init() {
	const defaultConfigDir = "./.config"
	cobra.OnInitialize(initOutput)

	rootCmd.PersistentFlags().StringVar(&configDir, "conf-dir", defaultConfigDir, "Directory of JWT clinet config")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format: json, text")

	rootCmd.AddCommand(create.GetCreateCommand())
}

func initOutput() {
	if err := config.InitConfig(configDir); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	logger.InitLogger(config.Get().EnableDebug, "")

	// TODO(set server addr)

	if err := output.Init(outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to set output option: %v\n", err)
		os.Exit(1)
	}
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
