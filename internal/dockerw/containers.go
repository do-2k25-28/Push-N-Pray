package dockerw

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/docker/go-sdk/container"
	tcontainer "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type ContainerNetwork struct {
	Name    string
	Aliases []string
}

type HealthConfig = tcontainer.HealthConfig
type RestartPolicy = tcontainer.RestartPolicy

type ContainerConfig struct {
	Image    string
	Name     string
	Networks []ContainerNetwork
	Env      map[string]string
	Labels   map[string]string
	Cmd      []string
	// Healthcheck uses Docker's native HEALTHCHECK support.
	Healthcheck *tcontainer.HealthConfig
	// RestartPolicy defaults to unless-stopped when unset.
	RestartPolicy *tcontainer.RestartPolicy
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

	if len(cfg.Env) > 0 {
		opts = append(opts, container.WithEnv(cfg.Env))
	}

	if len(cfg.Labels) > 0 {
		opts = append(opts, container.WithLabels(cfg.Labels))
	}

	if len(cfg.Cmd) > 0 {
		opts = append(opts, container.WithCmd(cfg.Cmd...))
	}

	opts = append(opts, container.WithAdditionalHostConfigModifier(func(hostConfig *tcontainer.HostConfig) {
		hostConfig.RestartPolicy = restartPolicyOrDefault(cfg.RestartPolicy)
	}))

	if cfg.Healthcheck != nil {
		opts = append(opts, container.WithAdditionalConfigModifier(func(config *tcontainer.Config) {
			config.Healthcheck = cfg.Healthcheck
		}))
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

func restartPolicyOrDefault(policy *tcontainer.RestartPolicy) tcontainer.RestartPolicy {
	if policy != nil {
		return *policy
	}

	return tcontainer.RestartPolicy{Name: tcontainer.RestartPolicyUnlessStopped}
}

// ExecInContainer runs a command inside a running container and returns an error if the exit code is non-zero.
func (c *Client) ExecInContainer(ctx context.Context, containerName string, cmd []string) error {
	args := append([]string{"exec", containerName}, cmd...)
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("exec in %s: %w: %s", containerName, err, out)
	}
	return nil
}

func NewHTTPHealthcheck(path string, port int, interval string, timeout string) (*tcontainer.HealthConfig, error) {
	parsedInterval, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("parse healthcheck interval: %w", err)
	}

	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("parse healthcheck timeout: %w", err)
	}

	return &tcontainer.HealthConfig{
		Test:     []string{"CMD-SHELL", fmt.Sprintf("wget --no-verbose --tries=1 --spider http://127.0.0.1:%d%s || curl --fail --silent http://127.0.0.1:%d%s >/dev/null", port, path, port, path)},
		Interval: parsedInterval,
		Timeout:  parsedTimeout,
	}, nil
}

func (c *Client) RunContainerFromConfig(ctx context.Context, config ContainerConfig) error {
	if _, err := container.Run(ctx, c.containerOptions(config)...); err != nil {
		return err
	}
	for _, net := range config.Networks {
		if err := c.ConnectContainerToNetwork(ctx, config.Name, net.Name); err != nil {
			return err
		}
	}
	return nil
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
