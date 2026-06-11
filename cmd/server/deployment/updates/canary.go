package updates

import (
	"context"
	"log"
	"pushnpray/internal/dockerw"
	"sort"
	"strconv"
	"strings"
	"time"
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

	sort.Slice(existingContainers, func(i, j int) bool {
		return existingContainers[i].Created > existingContainers[j].Created
	})

	router := newContainer.Name
	service := newContainer.Name
	weightedService := newContainer.Name + "-canary"
	oldService := strings.TrimPrefix(existingContainers[0].Names[0], "/")
	priority, _ := strconv.Atoi(existingContainers[0].Labels["traefik.http.routers."+oldService+".priority"])
	if priority == 0 {
		priority = 100
	}

	newContainer.Labels["traefik.http.routers."+router+".priority"] = strconv.Itoa(priority + 1)
	newContainer.Labels["traefik.http.routers."+router+".service"] = weightedService
	newContainer.Labels["traefik.http.services."+service+".loadbalancer.server.port"] = "80"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+service+".name"] = service + "@docker"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+service+".weight"] = "5"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+oldService+".name"] = oldService + "@docker"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+oldService+".weight"] = "95"

	log.Printf("Starting canary container %s with 5%% traffic\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	if err := docker.StopContainersByPattern(ctx, newContainer.Name); err != nil {
		return err
	}
	if err := docker.RemoveContainersByPattern(ctx, newContainer.Name); err != nil {
		return err
	}

	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+service+".weight"] = "25"
	newContainer.Labels["traefik.http.services."+weightedService+".weighted.services."+oldService+".weight"] = "75"

	log.Printf("Starting canary container %s with 25%% traffic\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	if err := docker.StopContainersByPattern(ctx, newContainer.Name); err != nil {
		return err
	}
	if err := docker.RemoveContainersByPattern(ctx, newContainer.Name); err != nil {
		return err
	}

	delete(newContainer.Labels, "traefik.http.routers."+router+".service")
	delete(newContainer.Labels, "traefik.http.services."+weightedService+".weighted.services."+service+".name")
	delete(newContainer.Labels, "traefik.http.services."+weightedService+".weighted.services."+service+".weight")
	delete(newContainer.Labels, "traefik.http.services."+weightedService+".weighted.services."+oldService+".name")
	delete(newContainer.Labels, "traefik.http.services."+weightedService+".weighted.services."+oldService+".weight")

	log.Printf("Starting canary container %s with 100%% traffic\n", newContainer.Name)
	if err := docker.RunContainerFromConfig(ctx, newContainer); err != nil {
		return err
	}

	log.Printf("Removing old containers matching %s\n", pattern)
	if err := docker.StopContainers(ctx, existingContainers); err != nil {
		return err
	}
	return docker.RemoveContainers(ctx, existingContainers)
}
