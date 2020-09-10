package project

import (
	"encoding/json"
	"io/ioutil"
	"os"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/output"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

var addProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Add New Project",
	Long:  "Add New Project",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		projectName, _ := cmd.Flags().GetString("name")

		if file == "" && projectName == "" {
			print.Error("\"name\" or \"file\" flag must be required")
			os.Exit(1)
		}

		if file != "" && projectName != "" {
			print.Error("Either \"id\" or \"file\" flag must be specified.")
			os.Exit(1)
		}

		token, err := config.GetAccessToken()
		if err != nil {
			print.Error("Token get failed: %v", err)
			os.Exit(1)
		}

		req := &projectapi.ProjectCreateRequest{}

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
			// set params by commandline arguments
			req.Name = projectName
			req.TokenConfig.AccessTokenLifeSpan, _ = cmd.Flags().GetUint("accessExpires")
			req.TokenConfig.RefreshTokenLifeSpan, _ = cmd.Flags().GetUint("refreshExpires")
			req.TokenConfig.SigningAlgorithm, _ = cmd.Flags().GetString("signAlg")
			req.AllowGrantTypes, _ = cmd.Flags().GetStringArray("grantTypes")

			pwPols, _ := cmd.Flags().GetStringArray("passwordPolicies")
			req.PasswordPolicy, err = util.ParsePolicies(pwPols)
			if err != nil {
				print.Error("Failed to parse password policy: %v", err)
				os.Exit(1)
			}

			req.UserLock.Enabled, _ = cmd.Flags().GetBool("userLockEnabled")
			if req.UserLock.Enabled {
				req.UserLock.MaxLoginFailure, _ = cmd.Flags().GetUint("maxLoginFailure")
				req.UserLock.LockDuration, _ = cmd.Flags().GetUint("lockDuration")
				req.UserLock.FailureResetTime, _ = cmd.Flags().GetUint("failureResetTime")
			}
		}

		c := config.Get()
		handler := apiclient.NewHandler(c.ServerAddr, token, c.Insecure, c.RequestTimeout)
		res, err := handler.ProjectAdd(req)
		if err != nil {
			print.Fatal("Failed to add project %s: %v", projectName, err)
		}

		format := output.NewProjectInfoFormat(res)
		output.Print(format)
	},
}

func init() {
	addProjectCmd.Flags().StringP("name", "n", "", "name of new project")
	addProjectCmd.Flags().Uint("accessExpires", 5*60, "access token life span [sec]")
	addProjectCmd.Flags().Uint("refreshExpires", 14*24*60*60, "refresh token life span [sec]")
	addProjectCmd.Flags().String("signAlg", "RS256", "token sigining algorithm, only support RS256")
	addProjectCmd.Flags().StringArray("grantTypes", []string{}, "allowed grant type list")
	addProjectCmd.Flags().StringArray("passwordPolicies", []string{}, "password policy of users, supports \"minLen=<uint>\", \"notUserName=<bool>\", \"useChar=<lower|upper|both|either>\", \"useDigit=<bool>\", \"useSpecialChar=<bool>\", \"blackLists=<string separated by semicolon(;)>\"")
	addProjectCmd.Flags().Bool("userLockEnabled", false, "enable user lock")
	addProjectCmd.Flags().Uint("maxLoginFailure", 5, "the max number of user login failure")
	addProjectCmd.Flags().Uint("lockDuration", 10*60, "a duration of couting login failure [sec]")
	addProjectCmd.Flags().Uint("failureResetTime", 10*60, "reset time of user locked [sec]")
	addProjectCmd.Flags().StringP("file", "f", "", "json file name of project info")
}
