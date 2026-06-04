package container

import (
	"context"
	"pushnpray/internal"
)

type Runner interface {
	RunContainerFromConfig(ctx context.Context, cfg internal.ContainerConfig) error
}
