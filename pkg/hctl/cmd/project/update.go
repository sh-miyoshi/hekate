package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

var updateProjectCmd = &cobra.Command{
	Use:   "update",
	Short: "Update project",
	Long:  "Update project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("name")
		file, _ := cmd.Flags().GetString("file")

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)

		req := &projectapi.ProjectPutRequest{}
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
			prev, err := handler.ProjectGet(projectName)
			if err != nil {
				print.Error("Failed to get previous client info: %v", err)
				os.Exit(1)
			}

			accessExpires := cmd.Flag("accessExpires")
			if accessExpires.Changed {
				a, _ := strconv.ParseUint(accessExpires.Value.String(), 10, 64)
				req.TokenConfig.AccessTokenLifeSpan = uint(a)
			} else {
				req.TokenConfig.AccessTokenLifeSpan = prev.TokenConfig.AccessTokenLifeSpan
			}
			refreshExpires := cmd.Flag("refreshExpires")
			if refreshExpires.Changed {
				r, _ := strconv.ParseUint(refreshExpires.Value.String(), 10, 64)
				req.TokenConfig.RefreshTokenLifeSpan = uint(r)
			} else {
				req.TokenConfig.RefreshTokenLifeSpan = prev.TokenConfig.RefreshTokenLifeSpan
			}
			signAlg := cmd.Flag("signAlg")
			if signAlg.Changed {
				req.TokenConfig.SigningAlgorithm = signAlg.Value.String()
			} else {
				req.TokenConfig.SigningAlgorithm = prev.TokenConfig.SigningAlgorithm
			}
			grantTypes := cmd.Flag("grantTypes")
			if grantTypes.Changed {
				req.AllowGrantTypes, _ = cmd.Flags().GetStringArray("grantTypes")
			} else {
				req.AllowGrantTypes = prev.AllowGrantTypes
			}
			passwordPolicies := cmd.Flag("passwordPolicies")
			if passwordPolicies.Changed {
				pwPols, _ := cmd.Flags().GetStringArray("passwordPolicies")
				req.PasswordPolicy, err = util.ParsePolicies(pwPols)
				if err != nil {
					print.Error("Failed to parse password policy: %v", err)
					os.Exit(1)
				}
			} else {
				req.PasswordPolicy = prev.PasswordPolicy
			}
		}

		if err := handler.ProjectUpdate(projectName, req); err != nil {
			print.Fatal("Failed to update project %s: %v", projectName, err)
		}

		print.Print("Successfully updated")
	},
}

func init() {
	updateProjectCmd.Flags().StringP("name", "n", "", "name of update project")
	updateProjectCmd.Flags().Uint("accessExpires", 5*60, "access token life span [sec]")
	updateProjectCmd.Flags().Uint("refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	updateProjectCmd.Flags().String("signAlg", "RS256", "token sigining algorithm, only support RS256")
	updateProjectCmd.Flags().StringArray("grantTypes", []string{}, "allowed grant type list")
	updateProjectCmd.Flags().StringArray("passwordPolicies", []string{}, "password policy of users, supports \"minLen=<uint>\", \"notUserName=<bool>\", \"useChar=<lower|upper|both|either>\", \"useDigit=<bool>\", \"useSpecialChar=<bool>\", \"blackLists=<string separated by semicolon(;)>\"")
	updateProjectCmd.Flags().StringP("file", "f", "", "json file name of project info")

	updateProjectCmd.MarkFlagRequired("name")
}
