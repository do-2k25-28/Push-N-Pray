package dockerw

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/go-sdk/image"
	dockerclient "github.com/moby/moby/client"
)

const buildArgLabelPrefix = "pushnpray.build-arg."

// BuildImage builds a Docker image with the given tag from the given Dockerfile and context directory.
// dockerfilePath may be absolute; it is resolved relative to contextDir for the SDK.
// buildArgs are passed as Docker build arguments and also stamped as image labels with the
// "pushnpray.build-arg." prefix so they can be verified with VerifyImageBuildArgs after the build.
func (c *Client) BuildImage(ctx context.Context, tag, dockerfilePath, contextDir string, buildArgs map[string]string) error {
	if contextDir == "" {
		contextDir = "."
	}

	relDockerfile, err := filepath.Rel(contextDir, dockerfilePath)
	if err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage: resolve dockerfile path: %w", err)
	}

	dockerBuildArgs := make(map[string]*string, len(buildArgs))
	labels := make(map[string]string, len(buildArgs))
	for k, v := range buildArgs {
		val := v
		dockerBuildArgs[k] = &val
		labels[buildArgLabelPrefix+k] = v
	}

	contextArchive, err := image.ArchiveBuildContext(contextDir, relDockerfile)
	if err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage: archive build context: %w", err)
	}

	fmt.Printf("Building docker image %s from %s...\n", tag, dockerfilePath)
	opts := dockerclient.ImageBuildOptions{
		Dockerfile: relDockerfile,
		BuildArgs:  dockerBuildArgs,
		Labels:     labels,
	}
	if _, err := image.Build(ctx, contextArchive, tag, image.WithBuildOptions(opts), image.WithBuildClient(c)); err != nil {
		return fmt.Errorf("dockerwrapper: BuildImage %q: %w", tag, err)
	}

	return nil
}

// VerifyImageBuildArgs inspects the built image and confirms that every key/value pair in
// expectedArgs is present as a "pushnpray.build-arg.<key>" label.  This provides a fast,
// inspectable proof that the intended build arguments were actually passed during the build.
func (c *Client) VerifyImageBuildArgs(ctx context.Context, tag string, expectedArgs map[string]string) error {
	result, err := c.ImageInspect(ctx, tag)
	if err != nil {
		return fmt.Errorf("dockerwrapper: VerifyImageBuildArgs: inspect %q: %w", tag, err)
	}

	labels := result.Config.Labels
	for k, expected := range expectedArgs {
		labelKey := buildArgLabelPrefix + k
		got, ok := labels[labelKey]
		if !ok {
			return fmt.Errorf("dockerwrapper: VerifyImageBuildArgs: build arg %q not found in image labels", k)
		}
		if got != expected {
			return fmt.Errorf("dockerwrapper: VerifyImageBuildArgs: build arg %q: expected %q, got %q", k, expected, got)
		}
	}

	return nil
}
