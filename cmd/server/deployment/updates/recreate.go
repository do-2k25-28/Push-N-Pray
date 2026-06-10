package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
)

type RecreateUpdateStrategy struct{}

func (s RecreateUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	log.Printf("Stopping containers matching %s\n", pattern)
	if err := docker.StopContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	log.Printf("Removing containers matching %s\n", pattern)
	if err := docker.RemoveContainersByPattern(ctx, pattern); err != nil {
		return err
	}

	log.Printf("Starting container %s\n", newContainer.Name)
	return docker.RunContainerFromConfig(ctx, newContainer)
}
