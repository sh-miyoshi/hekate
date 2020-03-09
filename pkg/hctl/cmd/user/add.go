package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	userapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/user"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

var addUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New User",
	Long:  "Add new user into the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		file, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")

		if file == "" && name == "" {
			print.Error("\"name\" or \"file\" flag must be required.")
			os.Exit(1)
		}

		if file != "" && name != "" {
			print.Error("Either \"name\" or \"file\" flag must be specified.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		req := &userapi.UserCreateRequest{}
		if file != "" {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				print.Error("Failed to read file %s: %v", file, err)
				os.Exit(1)
			}
			if err := json.Unmarshal(bytes, req); err != nil {
				print.Error("Failed to parse input file to json: %v", err)
				os.Exit(1)
			}
		} else {
			password, _ := cmd.Flags().GetString("password")
			if password == "" {
				// read password from console
				fmt.Printf("Password: ")
				var err error
				password, err = util.ReadPasswordFromConsole()
				if err != nil {
					print.Fatal("Failed to read password: %v", err)
				}
			}
			req.Name = name
			req.Password = password
			req.CustomRoles, _ = cmd.Flags().GetStringSlice("customRoles")
			req.SystemRoles, _ = cmd.Flags().GetStringSlice("systemRoles")
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		res, err := handler.UserAdd(projectName, req)
		if err != nil {
			print.Fatal("Failed to add new user %s to %s: %v", req.Name, projectName, err)
		}

		format := output.NewUserInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	addUserCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	addUserCmd.Flags().StringP("file", "f", "", "file path for new user info")
	addUserCmd.Flags().StringP("name", "n", "", "name of new user")
	addUserCmd.Flags().StringP("password", "p", "", "password of new user")
	addUserCmd.Flags().StringSlice("customRoles", nil, "custom role list")
	addUserCmd.Flags().StringSlice("systemRoles", nil, "system role list")
	addUserCmd.MarkFlagRequired("project")
}
