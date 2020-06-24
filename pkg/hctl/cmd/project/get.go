package project

import (
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var getProjectCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Projects in the cluster",
	Long:  "Get Projects in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		if projectName != "" {
			res, err := handler.ProjectGet(projectName)
			if err != nil {
				print.Fatal("Failed to get project %s: %v", projectName, err)
			}

			format := output.NewProjectInfoFormat(res)
			output.Print(format)
			return
		}

		res, err := handler.ProjectGetList()
		if err != nil {
			print.Fatal("Failed to get project: %v", err)
		}

		format := output.NewProjectsInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	getProjectCmd.Flags().StringP("name", "n", "", "name of project")
}
