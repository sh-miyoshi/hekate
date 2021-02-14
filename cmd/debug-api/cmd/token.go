package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sh-miyoshi/hekate/pkg/hctl/login"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/spf13/cobra"
)

func init() {
	tokenCmd.Flags().StringP("user", "u", "admin", "user name")
	tokenCmd.Flags().StringP("password", "p", "password", "password")
	tokenCmd.Flags().String("project", "master", "project")
	tokenCmd.Flags().String("server", "http://localhost:18443", "server address")
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "save token",
	Long:  "save token",
	Run: func(cmd *cobra.Command, args []string) {
		// get params
		serverAddr, _ := cmd.Flags().GetString("server")
		projectName, _ := cmd.Flags().GetString("project")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")

		// get token
		tkn, err := login.DoWithPassword(user, password, login.Info{
			ServerAddr:  serverAddr,
			Insecure:    true,
			Timeout:     10,
			ProjectName: projectName,
			ClientID:    "portal",
		})
		if err != nil {
			print.Error("Failed login with password: %v", err)
			os.Exit(1)
		}

		// save access token and refresh token in tmp/token.json
		outputFile := "tmp/token.json"
		bytes, _ := json.MarshalIndent(tkn, "", "  ")
		ioutil.WriteFile(outputFile, bytes, os.ModePerm)
		print.Print("Successfully save token info to %s", outputFile)
	},
}
