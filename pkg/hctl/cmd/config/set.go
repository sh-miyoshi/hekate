package config

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var setConfigCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config of hctl command",
	Long:  "Set config of hctl command",
	Run: func(cmd *cobra.Command, args []string) {
		// get previous cofig
		conf := config.Get()

		// update config
		if server, _ := cmd.Flags().GetString("server"); server != "" {
			conf.ServerAddr = server
		}
		if project, _ := cmd.Flags().GetString("project"); project != "" {
			conf.DefaultProject = project
		}
		if clientID, _ := cmd.Flags().GetString("client-id"); clientID != "" {
			conf.ClientID = clientID
		}
		if clientSecret, _ := cmd.Flags().GetString("client-secret"); clientSecret != "" {
			conf.ClientSecret = clientSecret
		}
		if timeout, _ := cmd.Flags().GetUint("timeout"); timeout != 0 {
			conf.RequestTimeout = timeout
		}
		if cmd.Flags().Changed("insecure") {
			insecure, _ := cmd.Flags().GetBool("insecure")
			conf.Insecure = insecure
		}

		// set and save
		config.Set(conf)
		if err := config.SaveToFile(); err != nil {
			print.Error("Failed to save config: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	setConfigCmd.Flags().String("server", "", "The address of hekate server")
	setConfigCmd.Flags().String("project", "", "Default project for operation")
	setConfigCmd.Flags().String("client-id", "", "Client ID of CLI tool")
	setConfigCmd.Flags().String("client-secret", "", "Client secret of CLI tool")
	setConfigCmd.Flags().Uint("timeout", 0, "Request timeout [sec]")
	setConfigCmd.Flags().Bool("insecure", false, "Access to hekate server without tls verify")
}
