package cmd

import (
	"fmt"
	"pushnpray/cmd/cli/prerun"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

var branch string
var tag string
var commit string

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Trigger a deployment",
	Long:    "Start a deployment for a project and optionally target a specific tag or commit",
	PreRunE: prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		_branch := branch
		if tag != "" || commit != "" {
			_branch = ""
		}

		res, err := client.DeployProject(cmd.Context(), man.ProjectId, api.DeployProjectRequest{
			Branch: _branch,
			Tag:    tag,
			Commit: commit,
		})

		if err != nil {
			return err
		}

		fmt.Printf("Deployment scheduled (id=%s)\n", res.ID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Branch to deploy")
	deployCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag to deploy")
	deployCmd.Flags().StringVarP(&commit, "commit", "c", "", "Commit to deploy")

	deployCmd.MarkFlagsMutuallyExclusive("branch", "tag", "commit")
}
