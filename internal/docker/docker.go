// Package dockerwrapper abstracts Docker operations for pulling images,
// creating networks, and creating containers.
package docker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	dockerSdk "github.com/docker/go-sdk/client"
	sdkcontainer "github.com/docker/go-sdk/container"
	sdkimage "github.com/docker/go-sdk/image"
	sdknetwork "github.com/docker/go-sdk/network"
	"github.com/moby/moby/api/types/container"
	dockerclient "github.com/moby/moby/client"
)

var defaultDockerSocket = "/var/run/docker/sock"

type Client struct {
	docker dockerSdk.SDKClient
}

// ContainerConfig holds the parameters for creating a container.
type ContainerConfig struct {
	Image          string
	Name           string
	Network        *sdknetwork.Network
	NetworkAliases []string
	Env            map[string]string
	Cmd            []string
	// ExposedPorts format: "8080/tcp".
	ExposedPorts []string
}

func NewClient(ctx context.Context) (*Client, error) {
	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = defaultDockerSocket
	}

	client, err := dockerSdk.New(ctx, dockerSdk.WithDockerHost(dockerHost))

	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: create client: %w", err)
	}

	return &Client{docker: client}, nil
}

// We only support docker hub so no registry
func registryCredentials(image string) (string, string, error) {
	return "", "", nil
}

// PullImages pulls images concurrently. All pulls are attempted; errors are
// collected and returned as a single joined error.
func (c *Client) PullImages(ctx context.Context, images ...string) error {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	for _, image := range images {
		wg.Add(1)
		go func(img string) {
			defer wg.Done()

			if err := sdkimage.Pull(ctx, img, sdkimage.WithPullClient(c.docker), sdkimage.WithCredentialsFn(registryCredentials)); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("pull %q: %w", img, err))
				mu.Unlock()
			}
		}(image)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("dockerwrapper: PullImages: %w", errors.Join(errs...))
	}

	return nil
}

func (c *Client) CreateNetwork(ctx context.Context, name string) (*sdknetwork.Network, error) {
	nw, err := sdknetwork.New(ctx,
		sdknetwork.WithName(name),
		sdknetwork.WithClient(c.docker),
	)

	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: CreateNetwork %q: %w", name, err)
	}

	return nw, nil
}

var ErrDockerStopRemoveFailed = errors.New("failed to stop and remove the container")

// CreateContainer creates a container without starting it.
func (c *Client) CreateContainer(ctx context.Context, cfg ContainerConfig) (*sdkcontainer.Container, error) {
	if cfg.Image == "" {
		return nil, fmt.Errorf("dockerwrapper: CreateContainer: Image is required")
	}

	if cfg.Name == "" {
		return nil, fmt.Errorf("dockerwrapper: CreateContainer: Name is required")
	}

	opts := []sdkcontainer.ContainerCustomizer{
		sdkcontainer.WithClient(c.docker),
		sdkcontainer.WithImage(cfg.Image),
		sdkcontainer.WithName(cfg.Name),
		sdkcontainer.WithNoStart(),
	}

	if cfg.Network != nil {
		opts = append(opts, sdkcontainer.WithNetwork(cfg.NetworkAliases, cfg.Network))
	}

	if len(cfg.Env) > 0 {
		opts = append(opts, sdkcontainer.WithEnv(cfg.Env))
	}

	if len(cfg.Cmd) > 0 {
		opts = append(opts, sdkcontainer.WithCmd(cfg.Cmd...))
	}

	if len(cfg.ExposedPorts) > 0 {
		opts = append(opts, sdkcontainer.WithExposedPorts(cfg.ExposedPorts...))
	}

	ctr, err := sdkcontainer.Run(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: CreateContainer %q: %w", cfg.Name, err)
	}

	return ctr, nil
}

// RunContainerFromConfig creates and starts a container using a full ContainerConfig.
func (c *Client) RunContainerFromConfig(ctx context.Context, cfg ContainerConfig) error {
	if cfg.Image == "" {
		return fmt.Errorf("dockerwrapper: RunContainerFromConfig: Image is required")
	}
	if cfg.Name == "" {
		return fmt.Errorf("dockerwrapper: RunContainerFromConfig: Name is required")
	}

	opts := []sdkcontainer.ContainerCustomizer{
		sdkcontainer.WithClient(c.docker),
		sdkcontainer.WithImage(cfg.Image),
		sdkcontainer.WithName(cfg.Name),
	}

	if cfg.Network != nil {
		opts = append(opts, sdkcontainer.WithNetwork(cfg.NetworkAliases, cfg.Network))
	}
	if len(cfg.Env) > 0 {
		opts = append(opts, sdkcontainer.WithEnv(cfg.Env))
	}
	if len(cfg.Cmd) > 0 {
		opts = append(opts, sdkcontainer.WithCmd(cfg.Cmd...))
	}
	if len(cfg.ExposedPorts) > 0 {
		opts = append(opts, sdkcontainer.WithExposedPorts(cfg.ExposedPorts...))
	}

	fmt.Printf("Running docker container %s from image %s...\n", cfg.Name, cfg.Image)
	if _, err := sdkcontainer.Run(ctx, opts...); err != nil {
		return fmt.Errorf("dockerwrapper: RunContainerFromConfig %q: %w", cfg.Name, err)
	}
	return nil
}

// BuildImage builds a Docker image with the given tag from the given Dockerfile and context directory.
// dockerfilePath may be absolute; it is resolved relative to contextDir for the SDK.
func (c *Client) BuildImage(ctx context.Context, tag, dockerfilePath, contextDir string) error {
	if contextDir == "" {
		contextDir = "."
	}
	relDockerfile, err := filepath.Rel(contextDir, dockerfilePath)
	if err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage: resolve dockerfile path: %w", err)
	}
	fmt.Printf("Building docker image %s from %s...\n", tag, dockerfilePath)
	if _, err := sdkimage.BuildFromDir(ctx, contextDir, relDockerfile, tag, sdkimage.WithBuildClient(c.docker)); err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage %q: %w", tag, err)
	}
	return nil
}

// RunContainer creates and starts a container with the given name from the given image.
func (c *Client) RunContainer(ctx context.Context, containerName, imageName string) error {
	fmt.Printf("Running docker container %s from image %s...\n", containerName, imageName)
	_, err := sdkcontainer.Run(ctx,
		sdkcontainer.WithClient(c.docker),
		sdkcontainer.WithName(containerName),
		sdkcontainer.WithImage(imageName),
	)
	if err != nil {
		return fmt.Errorf("dockerwrapper: RunContainer %q: %w", containerName, err)
	}
	return nil
}

func (c *Client) listContainersByPattern(ctx context.Context, pattern string) ([]container.Summary, error) {
	result, err := c.docker.ContainerList(ctx, dockerclient.ContainerListOptions{
		All:     true,
		Filters: make(dockerclient.Filters).Add("name", pattern),
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
		if _, err := c.docker.ContainerStop(ctx, ctr.ID, dockerclient.ContainerStopOptions{}); err != nil {
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
		if _, err := c.docker.ContainerRemove(ctx, ctr.ID, dockerclient.ContainerRemoveOptions{Force: true}); err != nil {
			errs = append(errs, fmt.Errorf("remove %s: %w", ctr.ID, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dockerwrapper: RemoveContainersByPattern: %w", errors.Join(errs...))
	}
	return nil
}
