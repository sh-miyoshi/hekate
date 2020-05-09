package client

import (
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var getClientCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Clients in the cluster",
	Long:  "Get Clients in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		clientID, _ := cmd.Flags().GetString("id")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		if clientID != "" {
			res, err := handler.ClientGet(projectName, clientID)
			if err != nil {
				print.Fatal("Failed to get client %s: %v", clientID, err)
			}

			format := output.NewClientInfoFormat(res)
			output.Print(format)
			return
		}

		res, err := handler.ClientGetList(projectName)
		if err != nil {
			print.Fatal("Failed to get client: %v", err)
		}

		format := output.NewClientsInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	getClientCmd.Flags().String("project", "", "[Required] name of the project to which the client belongs")
	getClientCmd.Flags().String("id", "", "id of the client")
	getClientCmd.MarkFlagRequired("project")
}
