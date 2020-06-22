package config

import "github.com/spf13/cobra"

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config of hctl command",
	Long:  "Get config of hctl command",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}
