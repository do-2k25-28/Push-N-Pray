package deployment

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
	"strings"
)

type DeployService struct{}

const traefikNet = "traefik"

func NewDeployService() *DeployService {
	return &DeployService{}
}

func (s *DeployService) DeployProject(projectSlug string, projectID string, m *manifest.Manifest, workspaceDir string) error {
	ctx := context.Background()
	dockerClient, err := internal.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	if _, err := dockerClient.CreateNetwork(ctx, traefikNet); err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create traefik network: %w", err)
	}

	for _, app := range m.Apps.Dockerfile {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
		imageName := fmt.Sprintf("%s-image", containerName)

		labels := traefikLabels(containerName, app.Name, projectSlug, projectID)

		dockerfilePath := utils.ResolvePath(workspaceDir, app.Dockerfile)
		contextPath := utils.ResolvePath(workspaceDir, app.Context)
		fmt.Printf("Deploying Dockerfile app: %s\n", app.Name)

		if err := dockerClient.BuildImage(ctx, imageName, dockerfilePath, contextPath); err != nil {
			return fmt.Errorf(errFmtAppDeployFailed+": %w", app.Name, err)
		}

		if err := dockerClient.RunContainerFromConfig(ctx, internal.ContainerConfig{
			Image:       imageName,
			Name:        containerName,
			NetworkName: traefikNet,
			Labels:      labels,
		}); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}

	for _, app := range m.Apps.Docker {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)

		labels := traefikLabels(containerName, app.Name, projectSlug, projectID)

		fmt.Printf("Deploying Docker image app: %s\n", app.Name)
		if err := dockerClient.RunContainerFromConfig(ctx, internal.ContainerConfig{
			Image:       app.Image,
			Name:        containerName,
			NetworkName: traefikNet,
			Labels:      labels,
		}); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}
	return nil
}

func traefikLabels(containerName, appName, projectSlug, projectID string) map[string]string {
	domain := fmt.Sprintf("%s-%s-%s.pushnpray.polydo.dev", appName, projectSlug, projectID)
	return map[string]string{
		"traefik.enable":         "true",
		"traefik.docker.network": traefikNet,
		fmt.Sprintf("traefik.http.routers.%s.rule", containerName):             fmt.Sprintf("Host(`%s`)", domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName):      "websecure",
		fmt.Sprintf("traefik.http.routers.%s.tls", containerName):              "true",
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", containerName): "le",
	}
}
