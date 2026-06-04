package docker

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/container"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
)

type ImageApp struct {
	container.App
}

func NewImageApp(app manifest.DockerApp, projectSlug, projectID string) ImageApp {
	return ImageApp{
		App: container.NewApp(app.Name, app.Image, projectSlug, projectID),
	}
}

func (app ImageApp) RunContainer(ctx context.Context, dockerClient *internal.Client) error {
	fmt.Printf("Deploying Docker image app: %s\n", app.Name())
	return app.Run(ctx, dockerClient)
}
