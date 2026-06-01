package cmd

import (
	"fmt"
	"pushnpray/internal/manifest"
	"pushnpray/internal/session"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Show project status",
	Long:         "Show information about the current project and the most recent deployments.",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := session.VerifyAuth(); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestPath, err := cmd.Root().Flags().GetString("file")
		if err != nil {
			manifestPath = "pushnpray.toml"
		}

		man, err := manifest.Unmarshal(manifestPath)
		if err != nil {
			return err
		}

		auth, err := session.GetAuthClientOption(man.Server)
		if err != nil {
			return err
		}

		client, err := api.NewClient(man.Server, auth)
		if err != nil {
			return err
		}

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
