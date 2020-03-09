package logout

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/logout"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
			print.Debug("Failed to get token: %v", err)
			return
		}

		if err := logout.Logout(cfg.ServerAddr, cfg.ProjectName, token); err != nil {
			print.Fatal("Logout failed: %v", err)
		}

		if err := config.RemoveSecretFile(); err != nil {
			print.Fatal("Failed to remove secret file: %v", err)
		}

		print.Print("Successfully logged out")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return logoutCmd
}
