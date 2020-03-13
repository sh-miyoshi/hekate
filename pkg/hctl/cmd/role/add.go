package role

import (
	"os"

	"github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	roleapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/customrole"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var addRoleCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New Role",
	Long:  "Add new role into the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		roleName, _ := cmd.Flags().GetString("name")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		req := &roleapi.CustomRoleCreateRequest{
			Name: roleName,
		}

		res, err := handler.RoleAdd(projectName, req)
		if err != nil {
			print.Fatal("Failed to add new role %s to %s: %v", req.Name, projectName, err)
		}

		format := output.NewCustomRoleFormat(res)
		output.Print(format)
	},
}

func init() {
	addRoleCmd.Flags().String("project", "", "[Required] name of the project to which the role belongs")
	addRoleCmd.Flags().StringP("name", "n", "", "[Required] name of new role")
	addRoleCmd.MarkFlagRequired("project")
	addRoleCmd.MarkFlagRequired("name")
}
