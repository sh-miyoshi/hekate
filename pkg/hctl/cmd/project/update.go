package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	apiclient "github.com/sh-miyoshi/hekate/pkg/apiclient/v1"
	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/v1/project"
	"github.com/sh-miyoshi/hekate/pkg/hctl/config"
	"github.com/sh-miyoshi/hekate/pkg/hctl/print"
	"github.com/sh-miyoshi/hekate/pkg/hctl/util"
	"github.com/spf13/cobra"
)

func getData(cmd *cobra.Command, name string, prev interface{}, typ string) interface{} {
	if !cmd.Flag(name).Changed {
		return prev
	}

	switch typ {
	case "string":
		return cmd.Flag(name).Value.String()
	case "uint":
		v, _ := strconv.ParseUint(cmd.Flag(name).Value.String(), 10, 64)
		return uint(v)
	case "stringarray":
		v, _ := cmd.Flags().GetStringArray(name)
		return v
	case "bool":
		v, _ := strconv.ParseBool(cmd.Flag(name).Value.String())
		return v
	case "time":
		v, _ := strconv.ParseUint(cmd.Flag(name).Value.String(), 10, 64)
		return time.Duration(v) * time.Second
	}

	return nil
}

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

			req.TokenConfig.AccessTokenLifeSpan = getData(cmd, "accessExpires", prev.TokenConfig.AccessTokenLifeSpan, "uint").(uint)
			req.TokenConfig.RefreshTokenLifeSpan = getData(cmd, "refreshExpires", prev.TokenConfig.RefreshTokenLifeSpan, "uint").(uint)
			req.TokenConfig.SigningAlgorithm = getData(cmd, "signAlg", prev.TokenConfig.SigningAlgorithm, "string").(string)
			req.AllowGrantTypes = getData(cmd, "grantTypes", prev.AllowGrantTypes, "stringarray").([]string)
			pwPols := getData(cmd, "passwordPolicies", prev.PasswordPolicy, "stringarray").([]string)
			req.PasswordPolicy, err = util.ParsePolicies(pwPols)
			if err != nil {
				print.Error("Failed to parse password policy: %v", err)
				os.Exit(1)
			}
			req.UserLock.Enabled = getData(cmd, "userLockEnabled", prev.UserLock.Enabled, "bool").(bool)
			req.UserLock.MaxLoginFailure = getData(cmd, "maxLoginFailure", prev.UserLock.MaxLoginFailure, "uint").(uint)
			req.UserLock.LockDuration = getData(cmd, "lockDuration", prev.UserLock.LockDuration, "time").(time.Duration)
			req.UserLock.FailureResetTime = getData(cmd, "failureResetTime", prev.UserLock.FailureResetTime, "time").(time.Duration)
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
	updateProjectCmd.Flags().Bool("userLockEnabled", false, "enable user lock")
	updateProjectCmd.Flags().Uint("maxLoginFailure", 5, "the max number of user login failure")
	updateProjectCmd.Flags().Uint("lockDuration", 10*60, "a duration of couting login failure [sec]")
	updateProjectCmd.Flags().Uint("failureResetTime", 10*60, "reset time of user locked [sec]")
	updateProjectCmd.Flags().StringP("file", "f", "", "json file name of project info")

	updateProjectCmd.MarkFlagRequired("name")
}
