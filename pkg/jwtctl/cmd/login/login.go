package login

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sh-miyoshi/jwt-server/pkg/jwtctl/config"
	"github.com/sh-miyoshi/jwt-server/pkg/logger"
	tokenapi "github.com/sh-miyoshi/jwt-server/pkg/tokenapi/v1"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
)

var (
	userName string
	password string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {
		serverAddr := config.Get().ServerAddr
		project := config.Get().ProjectName
		logger.Debug("server address: %s", serverAddr)

		if userName == "" {
			// Set user name from STDIN
			fmt.Printf("User Name: ")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			userName = stdin.Text()
		}

		if password == "" {
			// input password in STDIN
			fmt.Printf("Password: ")
			var err error
			password, err = readPasswordFromConsole()
			if err != nil {
				fmt.Printf("Failed to read password: %v\n", err)
				os.Exit(1)
			}
		}

		req := tokenapi.TokenRequest{
			Name:     userName,
			Secret:   password,
			AuthType: "password",
		}

		url := fmt.Sprintf("%s/api/v1/project/%s/token", serverAddr, project)
		body, err := json.Marshal(req)
		if err != nil {
			logger.Error("<Program Bug> Failed to marshal JSON: %v", err)
			os.Exit(1)
		}

		httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			logger.Error("<Program Bug> Failed to create http request: %v", err)
			os.Exit(1)
		}
		httpReq.Header.Add("Content-Type", "application/json")
		client := &http.Client{}
		httpRes, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("Failed to request server: %v", err)
			os.Exit(1)
		}
		defer httpRes.Body.Close()

		switch httpRes.StatusCode {
		case 200:
			var res tokenapi.TokenResponse
			if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
				logger.Error("<Program Bug> Failed to parse http response: %v", err)
				os.Exit(1)
			}
			// Output to secret file
			secretFile := filepath.Join(config.Get().ConfigDir, "secret")
			bytes, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				logger.Error("<Program Bug> Failed to marshal json: %v", err)
				os.Exit(1)
			}

			ioutil.WriteFile(secretFile, bytes, os.ModePerm)
			fmt.Println("Successfully logged in")
		case 401, 404:
			fmt.Println("Failed to login to system")
			fmt.Println("Please cheak user name or password (or project name)")
			os.Exit(1)
		case 500:
			fmt.Println("Internal Server Error is occured")
			fmt.Println("Please contact to your server administrator")
			os.Exit(1)
		default:
			logger.Error("<Program Bug> Unexpected http response code: %d", httpRes.StatusCode)
			os.Exit(1)
		}
	},
}

func init() {
	loginCmd.Flags().StringVarP(&userName, "name", "n", "", "Login User Name")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Login User Password")
}

// GetCommand ...
func GetCommand() *cobra.Command {
	return loginCmd
}

func readPasswordFromConsole() (string, error) {
	passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	password := string(passwordBytes)
	return password, nil
}
