package cmd

import (
	"os"

	"pushnpray/cmd/cli/cmd/env"
	"pushnpray/cmd/cli/cmd/pat"
	"pushnpray/cmd/cli/cmd/project"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pushnpray",
	Short: "Deploy and manage Push'n'Pray projects from your terminal",
	Long:  "The Push'n'Pray CLI lets you authenticate, initialize projects, trigger deployments, manage environment variables, and administer personal access tokens.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "pushnpray.toml", "Path to a project manifest")

	var _ = rootCmd.MarkFlagFilename("file", "toml")

	rootCmd.AddCommand(env.EnvCmd)
	rootCmd.AddCommand(pat.PatCmd)
	rootCmd.AddCommand(project.ProjectCmd)
}
