package deployment

import (
	"fmt"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal/manifest"
)

type DeployService struct{}

func NewDeployService() *DeployService {
	return &DeployService{}
}

func (s *DeployService) DeployProject(projectSlug string, projectID string, m *manifest.Manifest, workspaceDir string) error {
	for _, app := range m.Apps.Dockerfile {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
		imageName := fmt.Sprintf("%s-image", containerName)
		dockerfilePath := utils.ResolvePath(workspaceDir, app.Dockerfile)
		contextPath := utils.ResolvePath(workspaceDir, app.Context)
		fmt.Printf("Deploying Dockerfile app: %s\n", app.Name)

		if err := utils.BuildDockerImage(imageName, dockerfilePath, contextPath); err != nil {
			return fmt.Errorf(errFmtAppDeployFailed+": %w", app.Name, err)
		}
		if err := utils.RunDockerContainer(containerName, imageName); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}
	for _, app := range m.Apps.Docker {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
		fmt.Printf("Deploying Docker image app: %s\n", app.Name)
		if err := utils.RunDockerContainer(containerName, app.Image); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}
	return nil
}
