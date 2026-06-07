package dockerw

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/go-sdk/container"
	tcontainer "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type ContainerNetwork struct {
	Name    string
	Aliases []string
}

type ContainerConfig struct {
	Image    string
	Name     string
	Networks []ContainerNetwork
	Env      map[string]string
	Labels   map[string]string
	Cmd      []string
	// ExposedPorts format: "8080/tcp".
	ExposedPorts []string
	VolumeBinds  []string
}

func (c *Client) containerOptions(cfg ContainerConfig) []container.ContainerCustomizer {
	opts := []container.ContainerCustomizer{
		container.WithClient(c),
		container.WithImage(cfg.Image),
		container.WithName(cfg.Name),
	}

	for _, network := range cfg.Networks {
		opts = append(opts, container.WithNetworkName(network.Aliases, network.Name))
	}

	if len(cfg.Env) > 0 {
		opts = append(opts, container.WithEnv(cfg.Env))
	}

	if len(cfg.Labels) > 0 {
		opts = append(opts, container.WithLabels(cfg.Labels))
	}

	if len(cfg.Cmd) > 0 {
		opts = append(opts, container.WithCmd(cfg.Cmd...))
	}

	if len(cfg.ExposedPorts) > 0 {
		opts = append(opts, container.WithExposedPorts(cfg.ExposedPorts...))
	}

	if len(cfg.VolumeBinds) > 0 {
		opts = append(opts, container.WithAdditionalHostConfigModifier(func(hostConfig *tcontainer.HostConfig) {
			hostConfig.Binds = cfg.VolumeBinds
		}))
	}

	return opts
}

func (c *Client) RunContainerFromConfig(ctx context.Context, config ContainerConfig) error {
	_, err := container.Run(ctx, c.containerOptions(config)...)
	return err
}

func (c *Client) listContainersByPattern(ctx context.Context, pattern string) ([]tcontainer.Summary, error) {
	result, err := c.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: make(client.Filters).Add("name", pattern),
	})
	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: list containers: %w", err)
	}
	return result.Items, nil
}

// StopContainersByPattern stops all running containers whose names match the given pattern.
func (c *Client) StopContainersByPattern(ctx context.Context, pattern string) error {
	containers, err := c.listContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}
	var errs []error
	for _, ctr := range containers {
		if _, err := c.ContainerStop(ctx, ctr.ID, client.ContainerStopOptions{}); err != nil {
			errs = append(errs, fmt.Errorf("stop %s: %w", ctr.ID, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dockerwrapper: StopContainersByPattern: %w", errors.Join(errs...))
	}
	return nil
}

// RemoveContainersByPattern removes all containers whose names match the given pattern.
func (c *Client) RemoveContainersByPattern(ctx context.Context, pattern string) error {
	containers, err := c.listContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}
	var errs []error
	for _, ctr := range containers {
		if _, err := c.ContainerRemove(ctx, ctr.ID, client.ContainerRemoveOptions{Force: true}); err != nil {
			errs = append(errs, fmt.Errorf("remove %s: %w", ctr.ID, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dockerwrapper: RemoveContainersByPattern: %w", errors.Join(errs...))
	}
	return nil
}
