package dockerw

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func (c *Client) CreateVolumeIfNotExist(ctx context.Context, name string, labels map[string]string) error {
	if _, err := c.VolumeCreate(ctx, client.VolumeCreateOptions{Name: name, Labels: labels}); err != nil {
		return fmt.Errorf("create volume %q: %w", name, err)
	}

	return nil
}
