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

	for len(existingContainers) > 0 {
		container := existingContainers[0]
		log.Printf("Rolling out old container %s\n", container.ID)
		if err := docker.StopContainers(ctx, existingContainers[:1]); err != nil {
			return err
		}
		if err := docker.RemoveContainers(ctx, existingContainers[:1]); err != nil {
			return err
		}
		existingContainers = existingContainers[1:]
	}

	return nil
}
