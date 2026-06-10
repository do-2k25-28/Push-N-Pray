package cmd

import (
	"io"
	"os"
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var logsTail int
var logsFollow bool

var logsCmd = &cobra.Command{
	Use:     "logs <app>",
	Short:   "Fetch logs of an app",
	Long:    "Stream or print the last N lines of logs from a running app container",
	Args:    cobra.ExactArgs(1),
	PreRunE: prerun.Combine(prerun.Auth(), prerun.Manifest(), prerun.ApiClientFromManifest()),
	RunE: func(cmd *cobra.Command, args []string) error {
		man := prerun.GetManifest(cmd.Context())
		client := prerun.GetApiClient(cmd.Context())

		reader, err := client.GetAppLogs(cmd.Context(), man.ProjectId, args[0], logsTail, logsFollow)
		if err != nil {
			return err
		}

		defer func() { _ = reader.Close() }()

		_, err = io.Copy(os.Stdout, reader)
		return err
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVarP(&logsTail, "tail", "n", 100, "Number of lines to show from the end of the logs")
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
}
