package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	dockerclient "github.com/moby/moby/client"
)

type ServiceContainer struct {
	Name    string
	Labels  map[string]string
	Running bool
}

func (c *Client) EnsureNetwork(ctx context.Context, name string) error {
	_, err := c.docker.NetworkInspect(ctx, name, dockerclient.NetworkInspectOptions{})
	if err == nil {
		return nil
	}
	if !errdefs.IsNotFound(err) {
		return fmt.Errorf("inspect network %q: %w", name, err)
	}
	if _, err := c.docker.NetworkCreate(ctx, name, dockerclient.NetworkCreateOptions{}); err != nil {
		return fmt.Errorf("create network %q: %w", name, err)
	}
	return nil
}

func (c *Client) EnsureVolume(ctx context.Context, name string, labels map[string]string) error {
	if _, err := c.docker.VolumeCreate(ctx, dockerclient.VolumeCreateOptions{Name: name, Labels: labels}); err != nil {
		return fmt.Errorf("create volume %q: %w", name, err)
	}
	return nil
}

func (c *Client) ListServiceContainers(ctx context.Context, projectId string) ([]ServiceContainer, error) {
	result, err := c.docker.ContainerList(ctx, dockerclient.ContainerListOptions{
		All:     true,
		Filters: make(dockerclient.Filters).Add("label", "pushnpray.service=true").Add("label", "pushnpray.project-id="+projectId),
	})
	if err != nil {
		return nil, fmt.Errorf("list managed containers: %w", err)
	}

	serviceContainers := make([]ServiceContainer, 0, len(result.Items))
	for _, dockerContainer := range result.Items {
		serviceContainers = append(serviceContainers, ServiceContainer{
			Name:    strings.TrimPrefix(dockerContainer.Names[0], "/"),
			Labels:  dockerContainer.Labels,
			Running: dockerContainer.State == container.StateRunning,
		})
	}
	return serviceContainers, nil
}

func (c *Client) RemoveContainerAndVolume(ctx context.Context, containerName, volumeName string) error {
	_, _ = c.docker.ContainerStop(ctx, containerName, dockerclient.ContainerStopOptions{})
	_, containerErr := c.docker.ContainerRemove(ctx, containerName, dockerclient.ContainerRemoveOptions{Force: true})
	_, volumeErr := c.docker.VolumeRemove(ctx, volumeName, dockerclient.VolumeRemoveOptions{Force: true})
	return errors.Join(containerErr, volumeErr)
}

func (c *Client) StartContainer(ctx context.Context, name string) error {
	if _, err := c.docker.ContainerStart(ctx, name, dockerclient.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("start container %q: %w", name, err)
	}
	return nil
}
