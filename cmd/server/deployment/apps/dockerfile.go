package apps

import (
	"context"
	"fmt"
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

func (app DockerfileApp) LinkedApps() []string {
	return app.Links
}

func (app DockerfileApp) GetAllowOriginFrom() string {
	return app.AllowOriginFrom
}

func (app DockerfileApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error {
	return docker.BuildImage(ctx, app.imageTag(manifest), app.Dockerfile, app.Context)
}

func (app DockerfileApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.imageTag(manifest),
	}
}

func NewDockerFileApp(manifest manifest.DockerFileApp, workspace string) (DockerfileApp, error) {
	app := DockerfileApp{
		manifest,
	}

	dockerfile, err := utils.ResolvePath(workspace, manifest.Dockerfile)
	if err != nil {
		return DockerfileApp{}, fmt.Errorf("invalid dockerfile path: %w", err)
	}

	buildContext, err := utils.ResolvePath(workspace, manifest.Context)
	if err != nil {
		return DockerfileApp{}, fmt.Errorf("invalid context path: %w", err)
	}

	app.Dockerfile = dockerfile
	app.Context = buildContext

	return app, nil
}
