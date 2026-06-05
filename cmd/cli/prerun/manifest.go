package prerun

import (
	"context"
	"pushnpray/internal/manifest"

	"github.com/spf13/cobra"
)

type manifestCtxKeyType int

const ManifestContextKey manifestCtxKeyType = 0

func Manifest() CobraPreRun {
	return func(cmd *cobra.Command, args []string) error {
		man, err := manifest.Unmarshal(manifest.DefaultManifestName)
		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), ManifestContextKey, man)
		cmd.SetContext(ctx)

		return nil
	}
}

func GetManifest(ctx context.Context) *manifest.Manifest {
	return ctx.Value(ManifestContextKey).(*manifest.Manifest)
}
