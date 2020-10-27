package secret

import (
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get project secret",
	Long:  "Get project secret",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		res, err := handler.ProjectKeysGet(projectName)
		if err != nil {
			print.Fatal("Failed to get project secret info: %v", err)
		}

		format := output.NewKeysFormat(res)
		output.Print(format)
	},
}

func init() {
	getCmd.Flags().StringP("name", "n", "", "[Required] name of project")
	getCmd.MarkFlagRequired("name")
}
