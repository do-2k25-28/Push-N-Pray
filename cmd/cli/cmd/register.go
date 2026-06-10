package cmd

import (
	"fmt"
	"pushnpray/internal/session"
	"pushnpray/internal/utils"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var regEmail string
var regPassword string
var regToken string
var regServerUrl string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new Push'n'Pray account",
	Long:  "Register a new account by providing your email, a password, and a registration token. On success, you are automatically logged in.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !utils.IsValidEmail(regEmail) {
			return fmt.Errorf("email must be valid")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(regServerUrl)
		if err != nil {
			return err
		}

		response, err := client.Register(cmd.Context(), api.RegisterRequest{
			Email:         regEmail,
			Password:      regPassword,
			RegisterToken: regToken,
		})
		if err != nil {
			return err
		}

		fmt.Println("Account created. Successfully logged in.")

		return session.SaveBearerSession(regServerUrl, response.AccessToken, response.RefreshToken)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(&regEmail, "email", "u", "", "Account email")
	registerCmd.Flags().StringVarP(&regPassword, "password", "p", "", "Account password")
	registerCmd.Flags().StringVar(&regToken, "register-token", "", "Registration token")
	registerCmd.Flags().StringVar(&regServerUrl, "server", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")

	var _ = registerCmd.MarkFlagRequired("email")
	var _ = registerCmd.MarkFlagRequired("password")
	var _ = registerCmd.MarkFlagRequired("register-token")
}
