package pat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a personal access token",
	Long:  "Generate a new token with a name and optional expiration and display it once.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pat create called")
	},
}

func init() {
	PatCmd.AddCommand(createCmd)
}
