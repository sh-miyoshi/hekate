package jwtctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tokenConfigID       string
	tokenConfigPassword string
)

var tokenCmd = &cobra.Command{
	Use:   "gen-token",
	Short: "generate JWT token",
	Long:  `generate JWT token`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate token")
	},
}

func init() {
	tokenCmd.Flags().StringVar(&tokenConfigID, "id", "", "id of user")
	tokenCmd.Flags().StringVarP(&tokenConfigPassword, "password", "p", "", "psassword of user")
	tokenCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(tokenCmd)
}
