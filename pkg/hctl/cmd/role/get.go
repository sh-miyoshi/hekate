package role

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var getRoleCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Roles in the cluster",
	Long:  "Get Roles in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		name, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		res, err := handler.RoleGetList(projectName, name)
		if err != nil {
			print.Fatal("Failed to get role: %v", err)
		}

		if name != "" {
			format := output.NewRolesInfoFormat(res)
			output.Print(format)
		} else {
			format := output.NewCustomRoleFormat(res[0])
			output.Print(format)
		}
	},
}

func init() {
	getRoleCmd.Flags().String("project", "", "[Required] name of the project to which the role belongs")
	getRoleCmd.Flags().String("name", "", "name of the role")
	getRoleCmd.MarkFlagRequired("project")
}
