package cmd

import (
	"os"
	"pushnpray/internal/manifest"
	"pushnpray/internal/session"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var manifestPath string = ""

var projectName string
var repositoryUrl string
var initServerURL string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  `Create a project on the platform using the current repository metadata and store the project id locally.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Root().Flags().GetString("file")
		if err != nil {
			return err
		}

		manifestPath = file

		return session.VerifyAuth()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		authOption, err := session.GetAuthClientOption(initServerURL)
		if err != nil {
			return err
		}

		client, err := api.NewClient(initServerURL, authOption)
		if err != nil {
			return err
		}

		projectResponse, err := client.CreateProject(cmd.Context(), api.CreateProjectRequest{
			Slug:          projectName,
			RepositoryURL: repositoryUrl,
		})
		if err != nil {
			return err
		}

		man := manifest.Manifest{
			ProjectId:     projectResponse.ID,
			RepositoryUrl: repositoryUrl,
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

	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name for your new project")
	initCmd.Flags().StringVarP(&repositoryUrl, "repository", "r", "", "URL of the repository (must be http(s))")
	initCmd.Flags().StringVarP(&initServerURL, "server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")

	var _ = initCmd.MarkFlagRequired("name")
	var _ = initCmd.MarkFlagRequired("repository")
	var _ = initCmd.MarkFlagRequired("server")
}
