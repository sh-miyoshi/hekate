package login

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/login"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")

		// TODO(support authorization code flow)

		if projectName != "" {
			config.SetProjectName(projectName)
		}

		if userName == "" {
			// Set user name from STDIN
			fmt.Printf("User Name: ")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			userName = stdin.Text()
		}

		if password == "" {
			// input password in STDIN
			fmt.Printf("Password: ")
			var err error
			password, err = util.ReadPasswordFromConsole()
			if err != nil {
				print.Fatal("Failed to read password: %v", err)
			}
		}

		res, err := login.Do(config.Get().ServerAddr, config.Get().ProjectName, userName, password)
		if err != nil {
			print.Fatal("Failed to login: %v", err)
		}

		config.SetSecret(userName, res)
		print.Print("Successfully logged in")
	},
}

func init() {
	loginCmd.Flags().String("project", "", "name of the project to which the user belongs")
	loginCmd.Flags().StringP("name", "n", "", "Login User Name")
	loginCmd.Flags().StringP("password", "p", "", "Login User Password")
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return loginCmd
}
