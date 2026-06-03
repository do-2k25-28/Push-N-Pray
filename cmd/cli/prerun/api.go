package prerun

import (
	"context"
	"pushnpray/internal/manifest"
	"pushnpray/internal/session"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

type apiCtxKeyType int

const ApiContextKey apiCtxKeyType = 1

func ApiClient(cmd *cobra.Command, args []string) error {
	man := cmd.Context().Value(ManifestContextKey).(*manifest.Manifest)

	auth, err := session.GetAuthClientOption(man.Server)
	if err != nil {
		return err
	}

	client, err := api.NewClient(man.Server, auth)
	if err != nil {
		return err
	}

	ctx := context.WithValue(cmd.Context(), ApiContextKey, client)
	cmd.SetContext(ctx)

	return nil
}
