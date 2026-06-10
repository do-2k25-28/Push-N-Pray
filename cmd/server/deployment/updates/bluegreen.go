package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
)

type BlueGreenUpdateStrategy struct{}

func (s BlueGreenUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	existingContainers, err := docker.ListContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}

	log.Printf("Starting blue-green container %s\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	log.Printf("Stopping blue containers matching %s\n", pattern)
	if err := docker.StopContainers(ctx, existingContainers); err != nil {
		return err
	}

	log.Printf("Removing blue containers matching %s\n", pattern)
	return docker.RemoveContainers(ctx, existingContainers)
}
