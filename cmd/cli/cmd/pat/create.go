package pat

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"
	"pushnpray/pkg/api"
	"time"

	"github.com/spf13/cobra"
)

var createName string
var createDays int

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a personal access token",
	Long:         "Generate a new token with a name and optional expiration and display it once.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.ApiClientFromArg("server")),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := prerun.GetApiClient(cmd.Context())

		req := api.CreatePATRequest{
			Name: createName,
		}

		if createDays > 0 {
			expires := time.Now().Add(time.Duration(createDays) * 24 * time.Hour).Unix()
			req.ExpiresAt = &expires
		}

		resp, err := client.CreatePAT(cmd.Context(), req)
		if err != nil {
			return err
		}

		fmt.Printf("Token ID: %s\n", resp.ID)
		fmt.Printf("Token: %s\n", resp.Token)
		fmt.Println("Please copy this token now. You won't be able to see it again!")

		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Name of the token")
	createCmd.Flags().IntVarP(&createDays, "days", "d", 0, "Expiration in days (0 for no expiration)")
	createCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")
	var _ = createCmd.MarkFlagRequired("name")

	PatCmd.AddCommand(createCmd)
}
