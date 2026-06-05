package docker

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DeployableApp interface {
	RunContainer(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error
}
