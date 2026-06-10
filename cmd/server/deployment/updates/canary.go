package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
)

type CanaryUpdateStrategy struct{}

func (s CanaryUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	log.Printf("Starting canary container %s\n", newContainer.Name)
	return docker.RunContainerFromConfig(ctx, newContainer)
}
