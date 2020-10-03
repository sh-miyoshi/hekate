package update

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update User",
	Long:  "Update user",
	Run: func(cmd *cobra.Command, args []string) {
		print.Print("TODO: implement this")
	},
}

func init() {
	updateUserCmd.AddCommand(unlockUserCmd)

	updateUserCmd.Flags().String("project", "", "name of the project to which the user belongs")
	updateUserCmd.Flags().StringP("file", "f", "", "file path for update user info")
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return updateUserCmd
}
