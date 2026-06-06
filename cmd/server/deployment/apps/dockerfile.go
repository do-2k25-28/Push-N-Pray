package apps

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DockerfileApp struct {
	manifest.DockerFileApp
}

func (app *DockerfileApp) imageTag(manifest manifest.Manifest) string {
	// TODO?: add deployment ID
	return "img-" + manifest.ProjectId + "-" + app.Name
}

func (app *DockerfileApp) AppName() string {
	return app.Name
}

func (app *DockerfileApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	// TODO!: prefix path with cloned path, right now this doesn't build anything
	return docker.BuildImage(ctx, app.imageTag(manifest), app.Dockerfile, app.Context)
}

func (app *DockerfileApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.imageTag(manifest),
	}
}
