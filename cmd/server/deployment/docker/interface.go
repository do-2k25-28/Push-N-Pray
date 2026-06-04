package docker

import (
	"context"
	"pushnpray/internal"
)

type DeployableApp interface {
	RunContainer(ctx context.Context, dockerClient *internal.Client) error
}
