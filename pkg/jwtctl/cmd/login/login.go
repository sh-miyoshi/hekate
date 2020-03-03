package login

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/login"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/util"
	"github.com/spf13/cobra"
)

var (
	userName string
	password string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO(support authorization code flow)

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
				fmt.Printf("Failed to read password: %v\n", err)
				os.Exit(1)
			}
		}

		res, err := login.Do(config.Get().ServerAddr, config.Get().ProjectName, userName, password)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		config.SetSecret(userName, res)
		fmt.Println("Successfully logged in")
	},
}

func init() {
	loginCmd.Flags().StringVarP(&userName, "name", "n", "", "Login User Name")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Login User Password")
}

// Command ...
func Command() *cobra.Command {
	return loginCmd
}
