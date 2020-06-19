package client

import (
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var deleteClientCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Client",
	Long:  "Delete client from the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		id, _ := cmd.Flags().GetString("id")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		if err := handler.ClientDelete(projectName, id); err != nil {
			print.Fatal("Failed to delete the client %s to %s: %v", id, projectName, err)
		}

		print.Print("Client %s successfully deleted", id)
	},
}

func init() {
	deleteClientCmd.Flags().String("project", "", "[Required] name of the project to which the client belongs")
	deleteClientCmd.Flags().String("id", "", "[Required] id of the client")
	deleteClientCmd.MarkFlagRequired("project")
	deleteClientCmd.MarkFlagRequired("id")
}
