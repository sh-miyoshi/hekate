package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createProjectCmd)
	createCmd.AddCommand(createUserCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file or from arguments",
	Long:  `Create a resource from a file or from arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("create command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return createCmd
}
