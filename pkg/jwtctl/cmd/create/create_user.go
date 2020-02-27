package create

import (
	"fmt"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	userapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/user"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/spf13/cobra"
)

var user userapi.UserCreateRequest

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create New User",
	Long:  "Create new user into the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)
		req := &user

		res, err := handler.UserAdd(projectName, req)
		if err != nil {
			fmt.Printf("Failed to add new user %s to %s: %v", user.Name, project.Name, err)
			os.Exit(1)
		}

		format := output.NewUserInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	createUserCmd.Flags().StringVarP(&user.Name, "name", "n", "", "[Required] name of new user")
	createUserCmd.Flags().StringVarP(&user.Password, "password", "p", "", "[Required] password of new user")
	createUserCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	// TODO(system roles)
	// TODO(custom roles)
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("password")
	createUserCmd.MarkFlagRequired("project")
}
