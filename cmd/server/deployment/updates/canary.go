package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
	"strconv"
	"strings"
)

type CanaryUpdateStrategy struct{}

func (s CanaryUpdateStrategy) UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error {
	existingContainers, err := docker.ListContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}

	if len(existingContainers) == 0 {
		log.Printf("Starting first container %s\n", newContainer.Name)
		return docker.RunContainerFromConfig(ctx, newContainer)
	}

	router := newContainer.Name
	service := newContainer.Name
	weightedService := newContainer.Name + "-canary"

	newContainer.Labels["traefik.http.routers."+router+".priority"] = strconv.Itoa(100 + len(existingContainers))
	newContainer.Labels["traefik.http.routers."+router+".service"] = weightedService
	newContainer.Labels["traefik.http.services."+service+".loadbalancer.server.port"] = "80"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+service+".name"] = service + "@docker"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+service+".weight"] = "5"

	oldWeight := 95 / len(existingContainers)
	oldWeightRemainder := 95 % len(existingContainers)
	for i, container := range existingContainers {
		oldService := strings.TrimPrefix(container.Names[0], "/")
		weight := oldWeight
		if i == 0 {
			weight += oldWeightRemainder
		}
		newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+oldService+".name"] = oldService + "@docker"
		newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+oldService+".weight"] = strconv.Itoa(weight)
	}

	log.Printf("Starting canary container %s with 5%% traffic\n", newContainer.Name)
	return docker.RunContainerFromConfig(ctx, newContainer)
}
