package cmd

import (
	"fmt"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var email string
var password string
var token string
var serverUrl string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with email/password or a personal access token",
	Long:  "Start a session by exchanging credentials for tokens and storing them locally for future commands.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(serverUrl)
		if err != nil {
			return err
		}

		if password != "" {
			fmt.Println("Loggin in with", email, password)

			response, err := client.Login(cmd.Context(), api.LoginRequest{
				Email:    email,
				Password: password,
			})

			if err != nil {
				return err
			}

			fmt.Println("Logged in", response.AccessToken, response.RefreshToken)

			// TODO: write access and refresh token to file
		} else {
			fmt.Println("Logging in with PAT", token)
			// TODO: write pat to file
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&email, "email", "u", "", "Account email")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Account password")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "Personnal access token")
	loginCmd.Flags().StringVarP(&serverUrl, "server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")

	loginCmd.MarkFlagsOneRequired("password", "token")
	loginCmd.MarkFlagsMutuallyExclusive("password", "token")
	_ = loginCmd.MarkFlagRequired("email")
}
