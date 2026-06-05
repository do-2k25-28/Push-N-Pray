package docker

import (
	"context"
	"fmt"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DockerfileApp struct {
	manifest.DockerFileApp
}

func (app DockerfileApp) RunContainer(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error {
	name := manifest.ProjectId + "-" + app.Name
	img := "img-" + name

	if err := client.BuildImage(ctx, img, app.Dockerfile, app.Context); err != nil {
		return fmt.Errorf("failed to deploy app %s: %w", name, err)
	}

	container := dockerw.ContainerConfig{
		Name:  "app-" + name,
		Image: img,
	}

	return client.RunContainerFromConfig(ctx, container)
}
