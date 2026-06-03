package prerun

import (
	"pushnpray/internal/session"

	"github.com/spf13/cobra"
)

func Auth(cmd *cobra.Command, args []string) error {
	if err := session.VerifyAuth(); err != nil {
		return err
	}

	return nil
}
