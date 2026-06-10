package env

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"
	"pushnpray/pkg/api"
	"strings"

	"github.com/spf13/cobra"
)

var secretFlag bool

var setCmd = &cobra.Command{
	Use:          "set KEY=VALUE [KEY=VALUE ...]",
	Short:        "Set project environment variables",
	Long:         "Add or update environment variables for all services in a project using KEY=VALUE pairs.\nUse --secret to mark all variables in the call as secrets so their values are never returned by the API.",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		vars := make([]api.EnvVar, 0, len(args))
		for _, arg := range args {
			idx := strings.IndexByte(arg, '=')
			if idx <= 0 {
				return fmt.Errorf("invalid format %q: expected KEY=VALUE", arg)
			}
			vars = append(vars, api.EnvVar{
				Name:   arg[:idx],
				Value:  arg[idx+1:],
				Secret: secretFlag,
			})
		}

		if err := client.SetProjectEnv(cmd.Context(), man.ProjectId, api.SetProjectEnvRequest{
			Variables: vars,
		}); err != nil {
			return err
		}

		fmt.Printf("Set %d environment variable(s).\n", len(vars))
		return nil
	},
}

func init() {
	setCmd.Flags().BoolVarP(&secretFlag, "secret", "s", false, "Mark all variables in this call as secrets (values will never be returned by the API)")
	EnvCmd.AddCommand(setCmd)
}
