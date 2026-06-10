package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
	"time"
)

type RollingUpdateStrategy struct{}

func (s RollingUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	existingContainers, err := docker.ListContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}

	log.Printf("Starting container %s\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	log.Printf("Stopping containers matching %s\n", pattern)
	if err := docker.StopContainers(ctx, existingContainers); err != nil {
		return err
	}

	log.Printf("Removing containers matching %s\n", pattern)
	return docker.RemoveContainers(ctx, existingContainers)
}
