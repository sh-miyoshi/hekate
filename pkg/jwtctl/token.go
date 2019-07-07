package jwtctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"

	tokenapiv1 "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

type tokenConfig struct {
	ID       string
	Password string
	Output   string
}

var tokenConf tokenConfig

var tokenCmd = &cobra.Command{
	Use:   "gen-token",
	Short: "generate JWT token",
	Long:  `generate JWT token`,
	Run: func(cmd *cobra.Command, args []string) {
		// curl localhost:8080/api/v1/token -X POST -d {id:id, password:password}
		if len(tokenConf.Password) == 0 {
			fmt.Print("Password: ")
			password, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("Failed to get password")
				return
			}
			fmt.Println() // Add newline
			tokenConf.Password = string(password)
		}

		req := tokenapiv1.TokenCreateRequest{
			ID:       tokenConf.ID,
			Password: tokenConf.Password,
		}

		reqRaw, err := json.Marshal(req)
		if err != nil {
			fmt.Printf("Failed to parse input id and password: %v\n", err)
			return
		}

		resRaw, err := http.Post(
			globalConfig.ServerAddr+"/api/v1/token",
			"application/json",
			bytes.NewBuffer(reqRaw),
		)
		if err != nil {
			fmt.Printf("Failed to access server %s\n", globalConfig.ServerAddr)
			return
		}

		if resRaw.StatusCode != http.StatusOK {
			fmt.Printf("Failed to create token: %s\n", resRaw.Status)
			return
		}

		var res tokenapiv1.TokenCreateResponse
		if err := json.NewDecoder(resRaw.Body).Decode(&res); err != nil {
			fmt.Printf("Failed to decode response body: %v", err)
			return
		}

		switch tokenConf.Output {
		case "text":
			fmt.Println(res.Token)
		case "json":
			fmt.Printf("%v\n", res)
		case "pretty":
			fmt.Printf("Token: %s\n", res.Token)
		default:
			fmt.Printf("missing output option: %s\n", tokenConf.Output)
			fmt.Printf("Token: %s\n", res.Token)
		}
	},
}

func init() {
	tokenCmd.Flags().StringVar(&tokenConf.ID, "id", "", "id of user")
	tokenCmd.Flags().StringVarP(&tokenConf.Password, "password", "p", "", "psassword of user")
	tokenCmd.Flags().StringVarP(&tokenConf.Output, "output", "o", "text", "output type[text, json, pretty]")
	tokenCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(tokenCmd)
}
