package project

import (
	"github.com/sh-miyoshi/hekate/pkg/hctl/cmd/project/secret"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(addProjectCmd)
	projectCmd.AddCommand(deleteProjectCmd)
	projectCmd.AddCommand(getProjectCmd)
	projectCmd.AddCommand(updateProjectCmd)
	projectCmd.AddCommand(secret.GetCommand())
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project in the cluster",
	Long:  `Manage project in the cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		print.Error("project command requires subcommand")
	},
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return projectCmd
}
