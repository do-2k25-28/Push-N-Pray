package dockerw

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/client"
)

func (c *Client) CreateNetworkIfNotExist(ctx context.Context, name string) error {
	_, err := c.NetworkInspect(ctx, name, client.NetworkInspectOptions{})
	if err == nil {
		return nil
	}

	if !errdefs.IsNotFound(err) {
		return fmt.Errorf("inspect network %q: %w", name, err)
	}

	if _, err := c.NetworkCreate(ctx, name, client.NetworkCreateOptions{}); err != nil {
		return fmt.Errorf("create network %q: %w", name, err)
	}

	return nil
}
