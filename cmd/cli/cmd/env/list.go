package env

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List project environment variables",
	Long:         "Show the environment variables currently configured for a project. Secret variable values are not shown.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		resp, err := client.GetProjectEnv(cmd.Context(), man.ProjectId)
		if err != nil {
			return err
		}

		if len(resp.Variables) == 0 {
			fmt.Println("No environment variables set.")
			return nil
		}

		fmt.Printf("%-40s %s\n", "NAME", "VALUE")
		for _, v := range resp.Variables {
			if v.Secret {
				fmt.Printf("%-40s ********\n", v.Name)
			} else {
				fmt.Printf("%-40s %s\n", v.Name, *v.Value)
			}
		}

		return nil
	},
}

func init() {
	EnvCmd.AddCommand(listCmd)
}
