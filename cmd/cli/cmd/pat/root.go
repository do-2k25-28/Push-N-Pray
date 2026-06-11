package pat

import (
	"pushnpray/cmd/cli/prerun"

	"github.com/spf13/cobra"
)

var PatCmd = &cobra.Command{
	Use:               "pat",
	Short:             "Manage personal access tokens",
	PersistentPreRunE: prerun.Combine(prerun.Auth(), prerun.ApiClientFromArg("server")),
}

func init() {
	PatCmd.PersistentFlags().String("server", "https://api.pushnpray.hagridshut.net/v1/", "Push'N'Pray instance url")
}
