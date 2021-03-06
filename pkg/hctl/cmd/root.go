package cmd

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/client"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/login"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/logout"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/project"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/role"
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/user"
	globalconfig "github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"

	"os"
)

var (
	outputFormat string
	configDir    string
	debugMode    bool
)

func init() {
	cobra.OnInitialize(initOutput)

	rootCmd.PersistentFlags().StringVar(&configDir, "conf-dir", "", "Directory of hctl config")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Output debug message")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Set output format: json, text")

	rootCmd.AddCommand(login.GetCommand())
	rootCmd.AddCommand(logout.GetCommand())
	rootCmd.AddCommand(project.GetCommand())
	rootCmd.AddCommand(user.GetCommand())
	rootCmd.AddCommand(client.GetCommand())
	rootCmd.AddCommand(role.GetCommand())
	rootCmd.AddCommand(config.GetCommand())
}

func initOutput() {
	if err := globalconfig.InitConfig(configDir); err != nil {
		print.Error("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}

	print.Init(debugMode)

	if err := output.Init(outputFormat); err != nil {
		print.Error("Failed to set output option: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hctl",
	Short: "hctl is a command line tool for hekate",
	Long:  "hctl is a command line tool for hekate",
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
