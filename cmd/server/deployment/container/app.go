package container

import (
	"context"
	"fmt"
	"pushnpray/internal"
)

type App struct {
	name          string
	containerName string
	imageName     string
	projectId     string
	labels        map[string]string
}

func NewApp(name, imageName, projectSlug, projectId string) App {
	containerName := fmt.Sprintf("%s-%s-%s", name, projectSlug, projectId)
	return App{
		name:          name,
		containerName: containerName,
		imageName:     imageName,
		projectId:     projectId,
		labels:        traefikLabels(containerName, name, projectSlug, projectId),
	}
}

func (app App) Name() string {
	return app.name
}

func (app App) ImageName() string {
	return app.imageName
}

func (app App) config() internal.ContainerConfig {
	return internal.ContainerConfig{
		Image: app.imageName,
		Name:  app.containerName,
		Networks: []internal.ContainerNetwork{
			{Name: traefikNet},
			{Name: ProjectNetworkName(app.projectId)},
		},
		Labels: app.labels,
	}
}

func (app App) Run(ctx context.Context, dockerClient *internal.Client) error {
	if err := dockerClient.RunContainer(ctx, app.config()); err != nil {
		return fmt.Errorf("failed to run app %s: %w", app.name, err)
	}

	return nil
}
