package project

import (
	"fmt"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Long:  "Permanently remove a project and shut down any active deployments.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project delete called")
	},
}

func init() {
	ProjectCmd.AddCommand(deleteCmd)
}
