package cmd

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Show project status",
	Long:         "Show information about the current project and the most recent deployments.",
	SilenceUsage: true,
	PreRunE:      prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		project, err := client.GetProject(cmd.Context(), man.ProjectId)
		if err != nil {
			return err
		}

		fmt.Printf("Project ID: %s\n", project.ID)
		fmt.Printf("Slug: %s\n", project.Slug)
		fmt.Printf("Repository: %s\n", project.RepositoryURL)
		fmt.Printf("Created At: %s\n", project.CreatedAt)
		fmt.Printf("Updated At: %s\n", project.UpdatedAt)

		deployments, err := client.ListDeployments(cmd.Context(), project.ID)
		if err != nil {
			return err
		}

		fmt.Println("\nRecent Deployments:")
		count := 0
		for _, dep := range deployments.Deployments {
			if count >= 5 {
				break
			}
			fmt.Printf("- %s | Status: %s | Date: %s\n", dep.ID, dep.Status, dep.CreatedAt)
			count++
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
