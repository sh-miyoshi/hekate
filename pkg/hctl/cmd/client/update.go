package client

import (
	"encoding/json"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	clientapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/client"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

var updateClientCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a client",
	Long:  "Update a client in the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		file, _ := cmd.Flags().GetString("file")
		id, _ := cmd.Flags().GetString("id")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		req := &clientapi.ClientPutRequest{}
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
			prev, err := handler.ClientGet(projectName, id)
			if err != nil {
				print.Error("Failed to get previous client info: %v", err)
				os.Exit(1)
			}

			secret := cmd.Flag("secret")
			if secret.Changed {
				req.Secret = secret.Value.String()
			} else {
				req.Secret = prev.Secret
			}

			accessType := cmd.Flag("accessType")
			if accessType.Changed {
				at := accessType.Value.String()
				if at != "public" && at != "confidential" {
					print.Error("Invalid client type %s was specified.", at)
					os.Exit(1)
				}
				if at == "confidential" && req.Secret == "" {
					print.Error("Please set client secret if access type is confidential")
					os.Exit(1)
				}
				req.AccessType = at
			} else {
				req.AccessType = prev.AccessType
			}

			callbacks := cmd.Flag("callbacks")
			if callbacks.Changed {
				req.AllowedCallbackURLs, _ = cmd.Flags().GetStringSlice("callbacks")
			} else {
				req.AllowedCallbackURLs = prev.AllowedCallbackURLs
			}
		}

		if err := handler.ClientUpdate(projectName, id, req); err != nil {
			print.Fatal("Failed to update client %s to %s: %v", id, projectName, err)
		}

		print.Print("Successfully updated")
	},
}

func init() {
	updateClientCmd.Flags().String("project", "", "[Required] name of the project to which the client belongs")
	updateClientCmd.Flags().StringP("file", "f", "", "file path for new client info")
	updateClientCmd.Flags().String("id", "", "id of new client")
	updateClientCmd.Flags().String("secret", "", "secret of new client")
	updateClientCmd.Flags().String("accessType", "confidential", "access type of client (public or confidential)")
	updateClientCmd.Flags().StringSlice("callbacks", nil, "list of allowed callback url")

	updateClientCmd.MarkFlagRequired("project")
	updateClientCmd.MarkFlagRequired("id")
}
