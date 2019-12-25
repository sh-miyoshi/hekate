package delete

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.AddCommand(deleteProjectCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource from a file or from arguments",
	Long:  `Delete a resource from a file or from arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("delete command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return deleteCmd
}
