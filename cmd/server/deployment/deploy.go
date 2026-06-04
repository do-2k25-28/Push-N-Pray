package deployment

import (
	"context"
	"errors"
	"fmt"
	"pushnpray/cmd/server/deployment/container"
	"pushnpray/cmd/server/deployment/docker"
	"pushnpray/cmd/server/deployment/service"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
)

func DeployProject(projectSlug string, projectID string, projectManifest *manifest.Manifest, workspaceDir string) error {
	ctx := context.Background()
	dockerClient, err := internal.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	if err := dockerClient.EnsureNetwork(ctx, container.ProjectNetworkName(projectID)); err != nil {
		return fmt.Errorf("failed to create project network: %w", err)
	}

	serviceDefinitions, err := service.ServiceDefinitionsFromManifest(projectManifest)
	if err != nil {
		return fmt.Errorf("invalid services configuration: %w", err)
	}
	if err := service.UpdateServices(ctx, dockerClient, projectID, serviceDefinitions); err != nil {
		return fmt.Errorf("failed to update services: %w", err)
	}

	apps := make([]docker.DeployableApp, 0, len(projectManifest.Apps.Dockerfile)+len(projectManifest.Apps.Docker))

	for _, app := range projectManifest.Apps.Dockerfile {
		apps = append(apps, docker.NewDockerfileApp(app, projectSlug, projectID, workspaceDir))
	}

	for _, app := range projectManifest.Apps.Docker {
		apps = append(apps, docker.NewImageApp(app, projectSlug, projectID))
	}

	return runApps(ctx, dockerClient, apps)
}

func runApps(ctx context.Context, dockerClient *internal.Client, apps []docker.DeployableApp) error {
	var deployErrors []error
	for _, app := range apps {
		if err := app.RunContainer(ctx, dockerClient); err != nil {
			deployErrors = append(deployErrors, err)
		}
	}

	return errors.Join(deployErrors...)
}
