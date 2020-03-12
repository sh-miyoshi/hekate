package client

import (
	"github.com/spf13/cobra"
)

var deleteClientCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Client",
	Long:  "Delete client from the project",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
}
