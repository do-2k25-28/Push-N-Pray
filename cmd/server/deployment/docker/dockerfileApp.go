package docker

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/container"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
)

type DockerfileApp struct {
	container.App
	dockerfilePath string
	contextPath    string
}

func NewDockerfileApp(app manifest.DockerFileApp, projectSlug, projectID, workspaceDir string) DockerfileApp {
	containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
	return DockerfileApp{
		App:            container.NewApp(app.Name, containerName+"-image", projectSlug, projectID),
		dockerfilePath: utils.ResolvePath(workspaceDir, app.Dockerfile),
		contextPath:    utils.ResolvePath(workspaceDir, app.Context),
	}
}

func (app DockerfileApp) RunContainer(ctx context.Context, dockerClient *internal.Client) error {
	fmt.Printf("Deploying Dockerfile app: %s\n", app.Name())
	if err := dockerClient.BuildImage(ctx, app.ImageName(), app.dockerfilePath, app.contextPath); err != nil {
		return fmt.Errorf("failed to deploy app %s: %w", app.Name(), err)
	}
	return app.Run(ctx, dockerClient)
}
