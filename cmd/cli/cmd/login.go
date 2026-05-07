package cmd

import (
	"fmt"
	"pushnpray/internal/session"
	"pushnpray/internal/utils"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsValidEmail(email) {
			return fmt.Errorf("email must be valid")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if password != "" {
			client, err := api.NewClient(serverUrl)
			if err != nil {
				return err
			}

			response, err := client.Login(cmd.Context(), api.LoginRequest{
				Email:    email,
				Password: password,
			})

			if err != nil {
				return err
			}

			return session.SaveBearerSession(serverUrl, response.AccessToken, response.RefreshToken)
		}

		return session.SaveClassicSession(serverUrl, email, token)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&email, "email", "u", "", "Account email")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Account password")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "Personnal access token")
	loginCmd.Flags().StringVarP(&serverUrl, "server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")

	var _ = loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagsOneRequired("password", "token")
	loginCmd.MarkFlagsMutuallyExclusive("password", "token")
}
