package add

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addCmd.AddCommand(addProjectCmd)
	addCmd.AddCommand(addUserCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource from a file or from arguments",
	Long:  `Add a resource from a file or from arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("add command requires subcommand")
	},
}

// Command ...
func Command() *cobra.Command {
	return addCmd
}
