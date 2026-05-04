package pat

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Revoke a personal access token",
	Long:  "Delete a token by id so it can no longer be used for authentication.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pat delete called")
	},
}

func init() {
	PatCmd.AddCommand(deleteCmd)
}
