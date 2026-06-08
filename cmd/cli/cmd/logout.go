package cmd

import (
	"fmt"
	"pushnpray/internal/session"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear the current session",
	Long:  "Remove stored access and refresh tokens to prevent further authenticated calls.",
	PreRunE: prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		deleted, err := session.DeleteSession(serverUrl)
		if err != nil {
			return err
		}

		if !deleted {
			return fmt.Errorf("no session found for this server")
		}

		fmt.Println("Logged out.")
		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	logoutCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")
}
