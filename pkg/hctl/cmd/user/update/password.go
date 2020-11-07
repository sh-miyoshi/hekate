package update

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

var passwordChangeCmd = &cobra.Command{
	Use:   "password",
	Short: "Unlock User",
	Long:  "Unlock user",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		// input password in STDIN
		if password == "" {
			fmt.Printf("New Password: ")
			pw1, err := util.ReadPasswordFromConsole()
			if err != nil {
				print.Fatal("Failed to read password: %v", err)
			}
			fmt.Printf("Confirm: ")
			pw2, err := util.ReadPasswordFromConsole()
			if err != nil {
				print.Fatal("Failed to read password: %v", err)
			}

			if pw1 != pw2 {
				print.Fatal("Password and Confirmation are not same")
			}
			password = pw1
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		if err := handler.UserChangePassword(projectName, userName, password); err != nil {
			print.Fatal("Failed to change user %s password: %v", userName, err)
		}

		print.Print("User %s password successfully changed", userName)
	},
}

func init() {
	passwordChangeCmd.Flags().StringP("project", "p", "", "[Required] name of project to which the user belongs")
	passwordChangeCmd.Flags().StringP("name", "n", "", "[Required] name of target user")
	passwordChangeCmd.Flags().String("password", "", "new password")
	passwordChangeCmd.MarkFlagRequired("project")
	passwordChangeCmd.MarkFlagRequired("name")
}
