package dockerw

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func (c *Client) CreateNetworkIfNotExist(ctx context.Context, name string) error {
	out, err := exec.CommandContext(ctx, "docker", "network", "create", name).CombinedOutput()
	if err == nil || strings.Contains(string(out), "already exists") {
		return nil
	}
	return fmt.Errorf("create network %q: %s", name, out)
}

func (c *Client) ConnectContainerToNetwork(ctx context.Context, containerName, networkName string) error {
	out, err := exec.CommandContext(ctx, "docker", "network", "connect", networkName, containerName).CombinedOutput()
	if err == nil || strings.Contains(string(out), "already exists") {
		return nil
	}
	return fmt.Errorf("connect %s to %s: %s", containerName, networkName, out)
}
