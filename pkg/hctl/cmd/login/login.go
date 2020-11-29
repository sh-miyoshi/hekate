package login

import (
	oidcapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/oidc"
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
		clientID, _ := cmd.Flags().GetString("client-id")

		if projectName == "" {
			projectName = config.Get().DefaultProject
		}

		c := config.Get()
		var res *oidcapi.TokenResponse
		var err error

		if clientID != "" {
			// login with client credentials flow
			secret, _ := cmd.Flags().GetString("client-secret")
			if secret == "" {
				print.Print("Please set client-secret flag")
				return
			}
			res, err = login.DoWithClient(login.Info{
				ServerAddr:   c.ServerAddr,
				ProjectName:  projectName,
				ClientID:     clientID,
				ClientSecret: secret,
				Insecure:     c.Insecure,
				Timeout:      c.RequestTimeout,
			})
		} else {
			// login with device flow
			res, err = login.Do(login.Info{
				ServerAddr:   c.ServerAddr,
				ProjectName:  projectName,
				ClientID:     config.Get().ClientID,
				ClientSecret: config.Get().ClientSecret,
				Insecure:     c.Insecure,
				Timeout:      c.RequestTimeout,
			})
		}

		if err != nil {
			print.Fatal("Failed to login: %v", err)
		}

		config.SetSecret(projectName, res)
		print.Print("Successfully logged in")
	},
}

func init() {
	loginCmd.Flags().String("project", "", "name of the project to which the user belongs")
	loginCmd.Flags().String("client-id", "", "client id, if set this, client-secret flags also required")
	loginCmd.Flags().String("client-secret", "", "client secret")
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return loginCmd
}
