package logout

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/jwtctl/config"
	"github.com/sh-miyoshi/hekate/pkg/jwtctl/logout"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the system",
	Long:  `Logout from the system`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		token, err := config.GetRefreshToken()
		if err != nil {
			// TODO(Output err to debug message)
			return
		}

		if err := logout.Logout(cfg.ServerAddr, cfg.ProjectName, token); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := config.RemoveSecretFile(); err != nil {
			fmt.Printf("Failed to remove secret file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Successfully logged out")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return logoutCmd
}
