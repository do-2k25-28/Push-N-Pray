package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
)

type RollingUpdateStrategy struct{}

func (s RollingUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	log.Printf("Starting container %s\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	log.Printf("Stopping containers matching %s\n", pattern)
	if err := docker.StopContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	log.Printf("Removing containers matching %s\n", pattern)
	return docker.RemoveContainersByPattern(ctx, pattern)
}
