package client

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	clientapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/client"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
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
			print.Error("\"id\" or \"file\" flag must be required.")
			os.Exit(1)
		}

		if file != "" && id != "" {
			print.Error("Either \"id\" or \"file\" flag must be specified.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		req := &clientapi.ClientCreateRequest{}
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
			secret, _ := cmd.Flags().GetString("secret")
			accessType, _ := cmd.Flags().GetString("accessType")

			if accessType == "confidential" && secret == "" {
				// generate the secret internally
				secret = uuid.New().String()
			}

			req.ID = id
			req.Secret = secret
			req.AccessType = accessType
			req.AllowedCallbackURLs, _ = cmd.Flags().GetStringSlice("callbacks")
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		res, err := handler.ClientAdd(projectName, req)
		if err != nil {
			print.Fatal("Failed to add new client %s to %s: %v", req.ID, projectName, err)
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
