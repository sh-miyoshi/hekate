package config

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/spf13/cobra"
)

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config of hctl command",
	Long:  "Get config of hctl command",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Get()

		format := output.NewConfigFormat(conf)
		output.Print(format)
	},
}

func init() {
}
