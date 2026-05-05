package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with email/password or a personal access token",
	Long:  "Start a session by exchanging credentials for tokens and storing them locally for future commands.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Account username")
	loginCmd.Flags().StringP("password", "p", "", "Account password")
	loginCmd.Flags().StringP("token", "t", "", "Personnal access token")

	loginCmd.MarkFlagsOneRequired("password", "token")
	loginCmd.MarkFlagsMutuallyExclusive("password", "token")
	var _ = loginCmd.MarkFlagRequired("username")
}
