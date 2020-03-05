package user

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete User",
	Long:  "Delete User",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("user")

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		if err := handler.UserDelete(projectName, userName); err != nil {
			fmt.Printf("Failed to delete user %s: %v", userName, err)
			os.Exit(1)
		}

		fmt.Printf("User %s successfully deleted\n", userName)
	},
}

func init() {
	deleteUserCmd.Flags().StringP("project", "p", "", "[Required] name of project to which the user belongs")
	deleteUserCmd.Flags().StringP("name", "n", "", "[Required] name of delete user")
	deleteUserCmd.MarkFlagRequired("project")
	deleteUserCmd.MarkFlagRequired("name")
}
