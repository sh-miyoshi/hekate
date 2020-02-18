package get

import (
	"fmt"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/spf13/cobra"
)

var projectName string

var getProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Get Projects in the cluster",
	Long:  "Get Projects in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		// TODO(filtering by projectName)

		res, err := handler.ProjectGetList()
		if err != nil {
			fmt.Printf("Failed to get project: %v", err)
			os.Exit(1)
		}

		// TODO(output)
		for _, prj := range res {
			fmt.Printf("res: %v\n", prj)
		}

		// format := output.NewProjectInfoFormat(res)
		// output.Print(format)
	},
}

func init() {
	getProjectCmd.Flags().StringVarP(&projectName, "name", "n", "", "name of new project")
}
