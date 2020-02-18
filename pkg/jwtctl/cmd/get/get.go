package get

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	getCmd.AddCommand(getProjectCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources in the server",
	Long:  `Get resources in the server`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("get command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return getCmd
}
