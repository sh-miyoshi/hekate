package jwtctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jwtctl",
	Short: "jwtctl is a command tool for jwt-server",
	Long:  "jwtctl is a command tool for jwt-server",
	Run:   func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init(){
	// TODO(set server addr)
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
