package dockerw

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/go-sdk/image"
)

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
	if _, err := image.BuildFromDir(ctx, contextDir, relDockerfile, tag, image.WithBuildClient(c)); err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage %q: %w", tag, err)
	}

	return nil
}
