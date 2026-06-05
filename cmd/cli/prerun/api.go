package prerun

import (
	"context"
	"pushnpray/internal/manifest"
	"pushnpray/internal/session"
	"pushnpray/pkg/api"

	"github.com/spf13/cobra"
)

type apiCtxKeyType int
type serverCtxKeyType int

const ApiContextKey apiCtxKeyType = 1
const ServerContextKey serverCtxKeyType = 2

func ApiClientFromArg(arg string) CobraPreRun {
	return func(cmd *cobra.Command, args []string) error {
		server, err := cmd.Flags().GetString(arg)
		if err != nil {
			return err
		}

		return apiClient(server)(cmd, args)
	}
}

func ApiClientFromManifest() CobraPreRun {
	return func(cmd *cobra.Command, args []string) error {
		man := cmd.Context().Value(ManifestContextKey).(*manifest.Manifest)
		return apiClient(man.Server)(cmd, args)
	}
}

func apiClient(server string) CobraPreRun {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.WithValue(cmd.Context(), ServerContextKey, server)
		cmd.SetContext(ctx)

		auth, err := session.GetAuthClientOption(server)
		if err != nil {
			return err
		}

		client, err := api.NewClient(server, auth)
		if err != nil {
			return err
		}

		ctx2 := context.WithValue(cmd.Context(), ApiContextKey, client)
		cmd.SetContext(ctx2)

		return nil
	}
}

func GetApiClient(ctx context.Context) *api.Client {
	return ctx.Value(ApiContextKey).(*api.Client)
}

func GetServer(ctx context.Context) string {
	return ctx.Value(ServerContextKey).(string)
}
