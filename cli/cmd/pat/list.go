package pat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List personal access tokens",
	Long:  "Show the personal access tokens associated with the current account.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pat list called")
	},
}

func init() {
	PatCmd.AddCommand(listCmd)
}
