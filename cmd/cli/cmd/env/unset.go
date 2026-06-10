package env

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:          "unset KEY [KEY ...]",
	Short:        "Unset project environment variables",
	Long:         "Remove one or more environment variables from a project by name.",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		for _, name := range args {
			if err := client.DeleteProjectEnv(cmd.Context(), man.ProjectId, name); err != nil {
				return fmt.Errorf("failed to unset %q: %w", name, err)
			}
		}

		fmt.Printf("Unset %d environment variable(s).\n", len(args))
		return nil
	},
}

func init() {
	EnvCmd.AddCommand(unsetCmd)
}
