package jwtctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type config struct {
	ServerAddr string
}

var globalConfig config

var rootCmd = &cobra.Command{
	Use:   "jwtctl",
	Short: "jwtctl is a command tool for jwt-server",
	Long:  "jwtctl is a command tool for jwt-server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globalConfig.ServerAddr, "server", "http://localhost:8080", "server address")
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
