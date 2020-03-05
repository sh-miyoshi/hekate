package project

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(addProjectCmd)
	projectCmd.AddCommand(deleteProjectCmd)
	projectCmd.AddCommand(getProjectCmd)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project in the cluster",
	Long:  `Manage project in the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("project command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return projectCmd
}
