package env

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List project environment variables",
	Long:  "Show the environment variables currently configured for a project.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("env list called")
	},
}

func init() {
	EnvCmd.AddCommand(listCmd)
}
