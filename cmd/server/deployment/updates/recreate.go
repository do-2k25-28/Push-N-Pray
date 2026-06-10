package updates

import (
	"context"
	"pushnpray/internal/dockerw"
)

type RecreateUpdateStrategy struct{}

func (s RecreateUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	if err := docker.StopContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	if err := docker.RemoveContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	return docker.RunContainerFromConfig(ctx, newContainer)
}
