package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with email/password or a personal access token",
	Long:  "Start a session by exchanging credentials for tokens and storing them locally for future commands.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		password, _ := cmd.Flags().GetString("password")
		token, _ := cmd.Flags().GetString("token")

		if password == "" && token == "" {
			return fmt.Errorf("login requires either a password or a token")
		}

		if password != "" && token != "" {
			return fmt.Errorf("use either password or token but not both")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "Account username")
	loginCmd.Flags().StringP("password", "p", "", "Account password")
	loginCmd.Flags().StringP("token", "t", "", "Personnal access token")

	if err := loginCmd.MarkFlagRequired("username"); err != nil {
		panic(err)
	}
}
