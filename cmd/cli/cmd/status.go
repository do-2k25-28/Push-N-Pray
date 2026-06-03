package cmd

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"
	"pushnpray/internal/manifest"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Show project status",
	Long:         "Show information about the current project and the most recent deployments.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth, prerun.Manifest),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := cmd.Context().Value(prerun.ManifestContextKey).(*manifest.Manifest)
		client := cmd.Context().Value(prerun.ApiContextKey).(*api.Client)

		project, err := client.GetProject(cmd.Context(), man.ProjectId)
		if err != nil {
			return err
		}

		fmt.Printf("Project ID: %s\n", project.ID)
		fmt.Printf("Slug: %s\n", project.Slug)
		fmt.Printf("Repository: %s\n", project.RepositoryURL)
		fmt.Printf("Created At: %s\n", project.CreatedAt)
		fmt.Printf("Updated At: %s\n", project.UpdatedAt)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
