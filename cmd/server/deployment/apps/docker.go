package apps

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"

	"github.com/docker/go-sdk/image"
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

func (app DockerApp) GetAllowOriginFrom() string {
	return app.AllowOriginFrom
}

// We only support docker hub so no registry
func registryCredentials(image string) (string, string, error) {
	return "", "", nil
}

func (app DockerApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	return image.Pull(ctx, app.Image, image.WithPullClient(docker), image.WithCredentialsFn(registryCredentials))
}

func (app DockerApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.Image,
	}
}

func NewDockerApp(manifest manifest.DockerApp) DockerApp {
	return DockerApp{manifest}
}
