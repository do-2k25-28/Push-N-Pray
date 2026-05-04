package env

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset project environment variables",
	Long:  "Remove one or more environment variables from a project by name.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("env unset called")
	},
}

func init() {
	EnvCmd.AddCommand(unsetCmd)
}
