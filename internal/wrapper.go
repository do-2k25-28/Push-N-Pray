// Package dockerwrapper abstracts Docker operations for pulling images,
// creating networks, and creating containers.
package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"

	dockerSdk "github.com/docker/go-sdk/client"
	sdkcontainer "github.com/docker/go-sdk/container"
	sdkimage "github.com/docker/go-sdk/image"
	sdknetwork "github.com/docker/go-sdk/network"
)

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
	client, err := dockerSdk.New(ctx)

	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: create client: %w", err)
	}

	return &Client{docker: client}, nil
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
			if err := sdkimage.Pull(ctx, img, sdkimage.WithPullClient(c.docker)); err != nil {
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
