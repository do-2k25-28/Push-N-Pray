package internal

import (
	"context"
	"errors"
	"fmt"

	sdkcontainer "github.com/docker/go-sdk/container"
	"github.com/moby/moby/api/types/container"
)

type ContainerNetwork struct {
	Name    string
	Aliases []string
}

type ContainerConfig struct {
	Image        string
	Name         string
	Networks     []ContainerNetwork
	Environment  map[string]string
	Labels       map[string]string
	Command      []string
	VolumeBinds  []string
	ExposedPorts []string
}

func (c *Client) RunContainer(ctx context.Context, config ContainerConfig) error {
	if config.Image == "" || config.Name == "" {
		return errors.New("container image and name are required")
	}

	options := []sdkcontainer.ContainerCustomizer{
		sdkcontainer.WithClient(c.docker),
		sdkcontainer.WithImage(config.Image),
		sdkcontainer.WithName(config.Name),
	}
	for _, network := range config.Networks {
		options = append(options, sdkcontainer.WithNetworkName(network.Aliases, network.Name))
	}
	if len(config.Environment) > 0 {
		options = append(options, sdkcontainer.WithEnv(config.Environment))
	}
	if len(config.Labels) > 0 {
		options = append(options, sdkcontainer.WithLabels(config.Labels))
	}
	if len(config.Command) > 0 {
		options = append(options, sdkcontainer.WithCmd(config.Command...))
	}
	if len(config.VolumeBinds) > 0 {
		options = append(options, sdkcontainer.WithAdditionalHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Binds = config.VolumeBinds
		}))
	}
	if len(config.ExposedPorts) > 0 {
		options = append(options, sdkcontainer.WithExposedPorts(config.ExposedPorts...))
	}

	fmt.Printf("Running docker container %s from image %s...\n", config.Name, config.Image)
	if _, err := sdkcontainer.Run(ctx, options...); err != nil {
		return fmt.Errorf("run container %q: %w", config.Name, err)
	}
	return nil
}
