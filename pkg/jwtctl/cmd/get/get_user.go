package get

import (
	"fmt"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get users in the project",
	Long:  "Get users in the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		userName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		res, err := handler.UserGetList(projectName, userName)
		if err != nil {
			fmt.Printf("Failed to get user: %v", err)
			os.Exit(1)
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
