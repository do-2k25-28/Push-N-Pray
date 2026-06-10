package apps

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DockerApp struct {
	manifest.DockerApp
}

func (app DockerApp) AppName() string {
	return app.Name
}

func (app DockerApp) LinkedApps() []string {
	return app.Links
}

func (app DockerApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	return nil
}

func (app DockerApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.Image,
	}
}

func NewDockerApp(manifest manifest.DockerApp) DockerApp {
	return DockerApp{manifest}
}
