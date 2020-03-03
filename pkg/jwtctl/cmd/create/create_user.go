package create

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	userapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/user"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/util"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create New User",
	Long:  "Create new user into the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		file, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")

		if file == "" && name == "" {
			fmt.Println("\"name\" or \"file\" flag must be required")
			os.Exit(1)
		}

		if file != "" && name != "" {
			fmt.Printf("either \"name\" or \"file\" flag must be specified.")
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
	createUserCmd.Flags().String("project", "", "[Required] name of the project to which the user belongs")
	createUserCmd.Flags().StringP("file", "f", "", "file path for create user info")
	createUserCmd.Flags().StringP("name", "n", "", "name of new user")
	createUserCmd.Flags().StringP("password", "p", "", "password of new user")
	createUserCmd.Flags().StringSlice("customRoles", nil, "custom role list")
	createUserCmd.Flags().StringSlice("systemRoles", nil, "system role list")
	createUserCmd.MarkFlagRequired("project")
}
