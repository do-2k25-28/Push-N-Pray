package pat

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:          "delete <token-id>",
	Short:        "Revoke a personal access token",
	Long:         "Delete a token by id so it can no longer be used for authentication.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.ApiClientFromArg("server")),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := prerun.GetApiClient(cmd.Context())
		tokenId := args[0]

		err := client.DeletePAT(cmd.Context(), tokenId)
		if err != nil {
			return err
		}

		fmt.Printf("Token %s successfully deleted.\n", tokenId)
		return nil
	},
}

func init() {
	PatCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")
}
