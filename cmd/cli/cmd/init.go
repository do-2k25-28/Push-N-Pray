package cmd

import (
	"os"
	"pushnpray/internal/manifest"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var manifestPath string = ""

var projectName string
var repositoryUrl string

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

		// TODO: Check if user is logged in

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: POST v1/projects

		// TODO: Get it from the response
		projectId := uuid.New()

		man := manifest.Manifest{
			ProjectId:     projectId.String(),
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

	var _ = initCmd.MarkFlagRequired("name")
	var _ = initCmd.MarkFlagRequired("repository")
}
