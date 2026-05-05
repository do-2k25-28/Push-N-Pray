package env

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set project environment variables",
	Long:  "Add or update environment variables for all services in a project using key=value pairs or a file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("env set called")
	},
}

func init() {
	EnvCmd.AddCommand(setCmd)
}
