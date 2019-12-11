package cmd

import (
	"fmt"
	"github.com/sh-miyoshi/jwt-server/cmd/jwtctl/output"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/config"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/create"
	"github.com/sh-miyoshi/jwt-server/pkg/cmd/login"
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

	rootCmd.PersistentFlags().StringVar(&configDir, "conf-dir", defaultConfigDir, "Directory of jwtctl config")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format: json, text")

	rootCmd.AddCommand(create.GetCreateCommand())
	rootCmd.AddCommand(login.GetLoginCommand())
}

func initOutput() {
	if err := config.InitConfig(configDir); err != nil {
		fmt.Printf("[ERROR] Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	logger.InitLogger(config.Get().EnableDebug, "")

	if err := output.Init(outputFormat); err != nil {
		fmt.Printf("[ERROR] Failed to set output option: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "jwtctl",
	Short: "jwtctl is a command line tool for jwt-server",
	Long:  "jwtctl is a command line tool for jwt-server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute method run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
}
