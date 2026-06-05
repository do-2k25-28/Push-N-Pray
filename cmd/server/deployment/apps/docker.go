package docker

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DockerApp struct {
	manifest.DockerApp
}

func (app *DockerApp) RunContainer(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error {
	container := dockerw.ContainerConfig{
		Image:    app.Image,
		Name:     "app-" + app.Name + "-" + manifest.ProjectId,
		Networks: []dockerw.ContainerNetwork{},
	}

	return client.RunContainerFromConfig(ctx, container)
}
