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
	labels        map[string]string
}

func NewApp(name, imageName, projectSlug, projectID string) App {
	containerName := fmt.Sprintf("%s-%s-%s", name, projectSlug, projectID)
	return App{
		name:          name,
		containerName: containerName,
		imageName:     imageName,
		labels:        traefikLabels(containerName, name, projectSlug, projectID),
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
		Image:       app.imageName,
		Name:        app.containerName,
		NetworkName: traefikNet,
		Labels:      app.labels,
	}
}

func (app App) Run(ctx context.Context, runner Runner) error {
	if err := runner.RunContainerFromConfig(ctx, app.config()); err != nil {
		return fmt.Errorf("failed to run app %s: %w", app.name, err)
	}

	return nil
}
