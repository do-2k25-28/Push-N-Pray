package internal

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	dockerSdk "github.com/docker/go-sdk/client"
	sdkimage "github.com/docker/go-sdk/image"
	"github.com/moby/moby/api/types/container"
	dockerclient "github.com/moby/moby/client"
)

type Client struct {
	docker dockerSdk.SDKClient
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := dockerSdk.New(ctx)

	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: create client: %w", err)
	}

	return &Client{docker: client}, nil
}

var ErrDockerStopRemoveFailed = errors.New("failed to stop and remove the container")

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

// CheckIfDockerInstalled returns true if the Docker CLI is available in the system PATH.
func CheckIfDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}
