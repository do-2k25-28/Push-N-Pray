package project

import (
	"fmt"
	"os"
	"pushnpray/cmd/cli/prerun"
	"pushnpray/internal/manifest"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a project",
	Long:    "Permanently remove a project and shut down any active deployments.",
	PreRunE: prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		if err := client.DeleteProject(cmd.Context(), man.ProjectId); err != nil {
			return err
		}

		fmt.Println("Project deleted")

		if err := os.Remove(manifest.DefaultManifestName); err != nil {
			return err
		}

		fmt.Println("Manifest deleted")

		return nil
	},
}

func init() {
	ProjectCmd.AddCommand(deleteCmd)
}
