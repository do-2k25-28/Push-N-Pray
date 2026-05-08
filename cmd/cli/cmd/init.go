package cmd

import (
	"os"
	"pushnpray/internal/manifest"
	"pushnpray/internal/session"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  `Create a project on the platform using the current repository metadata and store the project id locally.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return session.VerifyAuth()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		repositoryURL, err := cmd.Flags().GetString("repository")
		if err != nil {
			return err
		}

		serverURL, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}

		manifestPath, err := cmd.Root().Flags().GetString("file")
		if err != nil {
			return err
		}

		authOption, err := session.GetAuthClientOption(serverURL)
		if err != nil {
			return err
		}

		client, err := api.NewClient(serverURL, authOption)
		if err != nil {
			return err
		}

		projectResponse, err := client.CreateProject(cmd.Context(), api.CreateProjectRequest{
			Slug:          projectName,
			RepositoryURL: repositoryURL,
		})
		if err != nil {
			return err
		}

		man := manifest.Manifest{
			ProjectId:     projectResponse.ID,
			RepositoryUrl: repositoryURL,
		}

		data, err := manifest.Marshal(&man)
		if err != nil {
			return err
		}

		return os.WriteFile(manifestPath, data, 0644)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("name", "n", "", "Name for your new project")
	initCmd.Flags().StringP("repository", "r", "", "URL of the repository (must be http(s))")
	initCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")

	var _ = initCmd.MarkFlagRequired("name")
	var _ = initCmd.MarkFlagRequired("repository")
}
