package user

import (
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:   "get",
	Short: "Get users in the project",
	Long:  "Get users in the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		res, err := handler.UserGetList(projectName, userName)
		if err != nil {
			print.Fatal("Failed to get user: %v", err)
		}

		format := output.NewUsersInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	getUserCmd.Flags().StringP("project", "p", "", "[Required] name of project to which the user belongs")
	getUserCmd.Flags().StringP("name", "n", "", "specific name user")
	getUserCmd.MarkFlagRequired("project")
}
