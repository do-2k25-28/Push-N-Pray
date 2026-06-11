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

func (app DockerfileApp) LinkedApps() []string {
	return app.Links
}

func (app DockerfileApp) GetAllowOriginFrom() string {
	return app.AllowOriginFrom
}

func (app DockerfileApp) Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest, env map[string]string) error {
	if err := docker.BuildImage(ctx, app.imageTag(manifest), app.Dockerfile, app.Context, env); err != nil {
		return err
	}
	if len(env) > 0 {
		return docker.VerifyImageBuildArgs(ctx, app.imageTag(manifest), env)
	}
	return nil
}

func (app DockerfileApp) ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig {
	return dockerw.ContainerConfig{
		Image: app.imageTag(manifest),
	}
}

func NewDockerFileApp(manifest manifest.DockerFileApp, workspace string) DockerfileApp {
	app := DockerfileApp{
		manifest,
	}

	app.Dockerfile = utils.ResolvePath(workspace, manifest.Dockerfile)
	app.Context = utils.ResolvePath(workspace, manifest.Context)

	return app
}
