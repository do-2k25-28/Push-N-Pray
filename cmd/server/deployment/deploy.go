package deployment

import (
	"context"
	"errors"
	"fmt"
	"pushnpray/cmd/server/deployment/docker"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
)

type DeployService struct{}

func NewDeployService() *DeployService {
	return &DeployService{}
}

func (s *DeployService) DeployProject(projectSlug string, projectID string, m *manifest.Manifest, workspaceDir string) error {
	ctx := context.Background()
	dockerClient, err := internal.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	apps := make([]docker.DeployableApp, 0, len(m.Apps.Dockerfile)+len(m.Apps.Docker))

	for _, app := range m.Apps.Dockerfile {
		apps = append(apps, docker.NewDockerfileApp(app, projectSlug, projectID, workspaceDir))
	}

	for _, app := range m.Apps.Docker {
		apps = append(apps, docker.NewImageApp(app, projectSlug, projectID))
	}

	return runApps(ctx, dockerClient, apps)
}

func runApps(ctx context.Context, client docker.Client, apps []docker.DeployableApp) error {
	var deployErrors []error
	for _, app := range apps {
		if err := app.RunContainer(ctx, client); err != nil {
			deployErrors = append(deployErrors, err)
		}
	}

	return errors.Join(deployErrors...)
}
