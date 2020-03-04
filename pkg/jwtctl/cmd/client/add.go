package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/jwt-server/pkg/apiclient/v1"
	clientapi "github.com/sh-miyoshi/jwt-server/pkg/apihandler/v1/client"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/output"
	"github.com/spf13/cobra"
)

var addClientCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New Client",
	Long:  "Add new client into the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		file, _ := cmd.Flags().GetString("file")
		id, _ := cmd.Flags().GetString("id")

		if file == "" && id == "" {
			fmt.Println("\"id\" or \"file\" flag must be required")
			os.Exit(1)
		}

		if file != "" && id != "" {
			fmt.Println("either \"id\" or \"file\" flag must be specified.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}

		req := &clientapi.ClientCreateRequest{}
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
			secret, _ := cmd.Flags().GetString("secret")
			accessType, _ := cmd.Flags().GetString("accessType")

			if accessType == "confidential" && secret == "" {
				fmt.Println("Please set client secret if access type is confidential")
				os.Exit(1)
			}

			req.ID = id
			req.Secret = secret
			req.AccessType = accessType
			req.AllowedCallbackURLs, _ = cmd.Flags().GetStringSlice("callbacks")
		}

		handler := apiclient.NewHandler(config.Get().ServerAddr, token)

		res, err := handler.ClientAdd(projectName, req)
		if err != nil {
			fmt.Printf("Failed to add new client %s to %s: %v", req.ID, projectName, err)
			os.Exit(1)
		}

		format := output.NewClientInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	addClientCmd.Flags().String("project", "", "[Required] name of the project to which the client belongs")
	addClientCmd.Flags().StringP("file", "f", "", "file path for new client info")
	addClientCmd.Flags().String("id", "", "id of new client")
	addClientCmd.Flags().String("secret", "", "secret of new client")
	addClientCmd.Flags().String("accessType", "confidential", "access type of client (public or confidential)")
	addClientCmd.Flags().StringSlice("callbacks", nil, "list of allowed callback url")
	addClientCmd.MarkFlagRequired("project")
}
