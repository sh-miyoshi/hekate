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
			fmt.Println("\"name\" or \"file\" flag must be required")
			os.Exit(1)
		}

		if file != "" && name != "" {
			fmt.Println("either \"name\" or \"file\" flag must be specified.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		req := &userapi.UserCreateRequest{}
		if file != "" {
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Failed to read file %s: %v\n", file, err)
				os.Exit(1)
			}
			if err := json.Unmarshal(bytes, req); err != nil {
				fmt.Printf("Failed to parse input file to json: %v\n", err)
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
					fmt.Printf("Failed to read password: %v\n", err)
					os.Exit(1)
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
			fmt.Printf("Failed to add new user %s to %s: %v", req.Name, projectName, err)
			os.Exit(1)
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
