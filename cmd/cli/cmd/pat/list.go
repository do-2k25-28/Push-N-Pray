package pat

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List personal access tokens",
	Long:         "Show the personal access tokens associated with the current account.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.ApiClientFromArg("server")),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := prerun.GetApiClient(cmd.Context())

		resp, err := client.ListPATs(cmd.Context())
		if err != nil {
			return err
		}

		if len(resp.Tokens) == 0 {
			fmt.Println("No personal access tokens found.")
			return nil
		}

		fmt.Printf("%-36s %-30s %-20s\n", "ID", "Name", "Expires At")
		for _, token := range resp.Tokens {
			expires := token.ExpiresAt
			if expires == "" {
				expires = "Never"
			}
			fmt.Printf("%-36s %-30s %-20s\n", token.ID, token.Name, expires)
		}

		return nil
	},
}

func init() {
	PatCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")
}
