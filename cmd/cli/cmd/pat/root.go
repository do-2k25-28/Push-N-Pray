package pat

import "github.com/spf13/cobra"

var PatCmd = &cobra.Command{
	Use:   "pat",
	Short: "Manage personal access tokens",
}
