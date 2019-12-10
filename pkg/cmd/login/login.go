package login

import (
	"github.com/spf13/cobra"
)

var (
	userName    string
	password    string
	projectName string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to system",
	Long:  `Login to system`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	loginCmd.Flags().StringVarP(&userName, "name", "n", "", "Login User Name")
	loginCmd.Flags().StringVar(&password, "password", "", "Login User Password")
	loginCmd.Flags().StringVar(&projectName, "project", "", "Set your project")
	loginCmd.MarkFlagRequired("name")
}

// GetLoginCommand ...
func GetLoginCommand() *cobra.Command {
	return loginCmd
}
