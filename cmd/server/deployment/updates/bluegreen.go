package updates

import (
	"context"
	"pushnpray/internal/dockerw"
)

type BlueGreenUpdateStrategy struct{}

func (s BlueGreenUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	if err := docker.StopContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	return docker.RemoveContainersByPattern(ctx, pattern)
}
