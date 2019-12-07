package create

import (
	"fmt"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create New Project",
	Long:  "Create New Project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create new project")
	},
}
