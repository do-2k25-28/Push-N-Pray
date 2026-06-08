package project

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List projects",
	Long:    "List every project owned by the authenticated user.",
	PreRunE: prerun.Combine(prerun.Auth(), prerun.ApiClientFromArg("server")),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := prerun.GetApiClient(cmd.Context())

		projects, err := client.ListProjects(cmd.Context())
		if err != nil {
			return err
		}

		if len(projects.Projects) == 0 {
			fmt.Println("No projects found")
			return nil
		}

		for _, project := range projects.Projects {
			fmt.Printf("%s\t%s\t%s\n", project.ID, project.Slug, project.RepositoryURL)
		}

		return nil
	},
	SilenceUsage: true,
}

func init() {
	ProjectCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("server", "", "https://api.pushnpray.polydo.dev/v1/", "Push'N'Pray instance url")
}
