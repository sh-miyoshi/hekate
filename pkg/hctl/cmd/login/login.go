package login

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/login"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")

		// TODO(login with client credentials)

		if projectName == "" {
			projectName = config.Get().DefaultProject
		}

		c := config.Get()
		res, err := login.Do(login.Info{
			ServerAddr:   c.ServerAddr,
			ProjectName:  projectName,
			ClientID:     config.Get().ClientID,
			ClientSecret: config.Get().ClientSecret,
			Insecure:     c.Insecure,
			Timeout:      c.RequestTimeout,
		})
		if err != nil {
			print.Fatal("Failed to login: %v", err)
		}

		config.SetSecret(projectName, res)
		print.Print("Successfully logged in")
	},
}

func init() {
	loginCmd.Flags().String("project", "", "name of the project to which the user belongs")
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return loginCmd
}
