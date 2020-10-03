package update

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var unlockUserCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock User",
	Long:  "Unlock user",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		if err := handler.UserUnlock(projectName, userName); err != nil {
			print.Fatal("Failed to unlock user %s: %v", userName, err)
		}

		print.Print("User %s successfully unlocked", userName)
	},
}

func init() {
	unlockUserCmd.Flags().StringP("project", "p", "", "[Required] name of project to which the user belongs")
	unlockUserCmd.Flags().StringP("name", "n", "", "[Required] name of target user")
	unlockUserCmd.MarkFlagRequired("project")
	unlockUserCmd.MarkFlagRequired("name")
}
