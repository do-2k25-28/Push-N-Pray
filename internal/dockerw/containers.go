package dockerw

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/go-sdk/container"
	"github.com/moby/moby/api/pkg/stdcopy"
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
		container.WithCredentialsFn(func(string) (string, string, error) {
			return "", "", nil
		}),
		container.WithAdditionalHostConfigModifier(func(hostConfig *tcontainer.HostConfig) {
			hostConfig.RestartPolicy = tcontainer.RestartPolicy{Name: tcontainer.RestartPolicyUnlessStopped}
			hostConfig.PublishAllPorts = false
		}),
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

	if len(cfg.Networks) > 0 {
		for _, net := range cfg.Networks {
			opts = append(opts, container.WithNetworkName(net.Aliases, net.Name))
		}
	}

	if len(cfg.VolumeBinds) > 0 {
		opts = append(opts, container.WithAdditionalHostConfigModifier(func(hostConfig *tcontainer.HostConfig) {
			hostConfig.Binds = cfg.VolumeBinds
		}))
	}

	return opts
}

// ExecInContainer runs a command inside a running container and returns an error if the exit code is non-zero.
func (c *Client) ExecInContainer(ctx context.Context, containerName string, cmd []string) error {
	created, err := c.ExecCreate(ctx, containerName, client.ExecCreateOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return fmt.Errorf("exec in %s: %w", containerName, err)
	}

	attached, err := c.ExecAttach(ctx, created.ID, client.ExecAttachOptions{})
	if err != nil {
		return fmt.Errorf("exec in %s: attach: %w", containerName, err)
	}
	defer attached.Close()

	var out bytes.Buffer
	if _, err := stdcopy.StdCopy(&out, &out, attached.Reader); err != nil {
		return fmt.Errorf("exec in %s: read output: %w", containerName, err)
	}

	inspected, err := c.ExecInspect(ctx, created.ID, client.ExecInspectOptions{})
	if err != nil {
		return fmt.Errorf("exec in %s: inspect: %w", containerName, err)
	}
	if inspected.ExitCode != 0 {
		return fmt.Errorf("exec in %s: exit code %d: %s", containerName, inspected.ExitCode, out.String())
	}
	return nil
}

func (c *Client) RunContainerFromConfig(ctx context.Context, config ContainerConfig) error {
	if _, err := container.Run(ctx, c.containerOptions(config)...); err != nil {
		return err
	}

	return nil
}

func (c *Client) ListContainersByPattern(ctx context.Context, pattern string) ([]tcontainer.Summary, error) {
	result, err := c.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: make(client.Filters).Add("name", pattern),
	})
	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: list containers: %w", err)
	}
	return result.Items, nil
}

func (c *Client) StopContainers(ctx context.Context, containers []tcontainer.Summary) error {
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

// GetContainerLogs returns a reader for the logs of a named container.
// The caller must close the returned reader.
func (c *Client) GetContainerLogs(ctx context.Context, containerName, tail string, follow bool) (io.ReadCloser, error) {
	reader, err := c.ContainerLogs(ctx, containerName, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       tail,
		Follow:     follow,
	})
	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: get logs for %s: %w", containerName, err)
	}
	return reader, nil
}

// StopContainersByPattern stops all running containers whose names match the given pattern.
func (c *Client) StopContainersByPattern(ctx context.Context, pattern string) error {
	containers, err := c.ListContainersByPattern(ctx, pattern)

	if err != nil {
		return err
	}

	return c.StopContainers(ctx, containers)
}

func (c *Client) RemoveContainers(ctx context.Context, containers []tcontainer.Summary) error {
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

// RemoveContainersByPattern removes all containers whose names match the given pattern.
func (c *Client) RemoveContainersByPattern(ctx context.Context, pattern string) error {
	containers, err := c.ListContainersByPattern(ctx, pattern)
	if err != nil {
		return err
	}
	return c.RemoveContainers(ctx, containers)
}
