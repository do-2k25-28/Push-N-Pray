package apps

import (
	"context"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DockerfileApp struct {
	manifest.DockerFileApp
}

func (app DockerfileApp) imageTag(manifest manifest.Manifest) string {
	// TODO?: add deployment ID
	return "img-" + manifest.ProjectId + "-" + app.Name
}

func (app DockerfileApp) AppName() string {
	return app.Name
}

func (app DockerfileApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	return docker.BuildImage(ctx, app.imageTag(manifest), app.Dockerfile, app.Context)
}

func (app DockerfileApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) (dockerw.ContainerConfig, error) {
	healthcheck, err := healthcheckConfig(app.App)
	if err != nil {
		return dockerw.ContainerConfig{}, err
	}

	return dockerw.ContainerConfig{
		Image:       app.imageTag(manifest),
		Healthcheck: healthcheck,
	}, nil
}

func NewDockerFileApp(manifest manifest.DockerFileApp, workspace string) DockerfileApp {
	app := DockerfileApp{
		manifest,
	}

	app.Dockerfile = utils.ResolvePath(workspace, manifest.Dockerfile)
	app.Context = utils.ResolvePath(workspace, manifest.Context)

	return app
}
