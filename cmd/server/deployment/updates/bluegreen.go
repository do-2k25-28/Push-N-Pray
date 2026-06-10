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

	traefikNetwork := newContainer.Labels["traefik.docker.network"]
	networks := newContainer.Networks
	newContainer.Networks = nil
	for _, network := range networks {
		if network.Name != traefikNetwork {
			newContainer.Networks = append(newContainer.Networks, network)
		}
	}

	log.Printf("Starting green container %s without external traffic\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	log.Printf("Switching traffic to green container %s on network %s\n", newContainer.Name, traefikNetwork)
	if err := docker.ConnectContainerToNetwork(ctx, newContainer.Name, traefikNetwork); err != nil {
		return err
	}

	log.Printf("Stopping blue containers matching %s after switching traffic\n", pattern)
	if err := docker.StopContainers(ctx, existingContainers); err != nil {
		return err
	}

	log.Printf("Removing blue containers matching %s\n", pattern)
	if err := docker.RemoveContainers(ctx, existingContainers); err != nil {
		return err
	}

	return nil
}
