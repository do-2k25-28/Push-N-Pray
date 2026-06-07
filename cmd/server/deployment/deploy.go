package deployment

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/apps"
	"pushnpray/cmd/server/deployment/project"
	"pushnpray/cmd/server/deployment/services"
	s3infra "pushnpray/infrastructure/s3"
	"pushnpray/internal/dockerw"

	"pushnpray/cmd/server/deployment/project"
	"pushnpray/cmd/server/deployment/services"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
	"pushnpray/internal/utils"
)

type DeployService struct {
	cephEndpoint string
}

func NewDeployService(cephEndpoint string) *DeployService {
	return &DeployService{cephEndpoint: cephEndpoint}
}

func (s *DeployService) DeployProject(projectSlug string, projectID string, m *manifest.Manifest, workspaceDir string) error {
	ctx := context.Background()

	dockerClient, err := internal.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create project Docker network
	projectNetwork := project.NetworkName(projectID)
	if err := dockerClient.EnsureNetworkExists(ctx, projectNetwork); err != nil {
		return fmt.Errorf("failed to create project network: %w", err)
	}

	// Manage dockerfile apps in the manifest
	for _, app := range m.Apps.Dockerfile {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
		imageName := fmt.Sprintf("%s-image", containerName)
		dockerfilePath := utils.ResolvePath(workspaceDir, app.Dockerfile)
		contextPath := utils.ResolvePath(workspaceDir, app.Context)
		fmt.Printf("Deploying Dockerfile app: %s\n", app.Name)

	if err := docker.CreateNetworkIfNotExist(ctx, project.NetworkName(manifest.ProjectId)); err != nil {
		return fmt.Errorf("failed to create project network: %w", err)
	}

	network := dockerw.ContainerNetwork{
		Name:    project.NetworkName(manifest.ProjectId),
		Aliases: []string{},
	}
	traefikNet := dockerw.ContainerNetwork{
		Name:    project.TraefikNet(),
		Aliases: []string{},
	}

	// Handle managed services creation and deletion

	_services := make([]services.ManagedService, 0, manifest.GetServiceCount())
	envsFromServices := []map[string]map[string]string{}

	for _, service := range manifest.Services.Postgres {
		pg := services.PostgresService{Manifest: service}
		_services = append(_services, &pg)
	}

	for _, service := range _services {
		deployed, err := service.IsDeployed(ctx, manifest)
		if err != nil {
			return err
		}

		if !deployed {
			if err := service.Prepare(ctx, docker, manifest); err != nil {
				return fmt.Errorf("filed to prepare deployment of service")
			}

			if err := service.Deploy(ctx, docker, manifest, network); err != nil {
				return fmt.Errorf("failed to deploy service")
			}
		}

		env, err := service.EnvToInject(manifest)
		if err != nil {
			return err
		}
		envsFromServices = append(envsFromServices, env)
	}

	appToEnv := utils.MergeMaps(envsFromServices)

	// Deploy or update application containers

	_apps := make([]apps.DeployableApp, 0, manifest.GetApplicationCount())

	for _, app := range manifest.Apps.Docker {
		_apps = append(_apps, apps.NewDockerApp(app))
	}

	for _, app := range manifest.Apps.Dockerfile {
		_apps = append(_apps, apps.NewDockerFileApp(app, workspaceDir))
	}

	for _, app := range _apps {
		if err := app.Prepare(ctx, docker, manifest); err != nil {
			return err
		}

		config := app.ContainerConfig(ctx, manifest)

		config.Env = utils.MergeMap(
			config.Env,
			appToEnv[app.AppName()], // Override user defined vars if they overlap
		)

		config.Name = "app-" + manifest.ProjectId + "-" + app.AppName()
		config.Networks = []dockerw.ContainerNetwork{network, traefikNet}
		config.Labels = utils.MergeMap(
			config.Labels,
			project.TraefikLabels(config.Name, app.AppName(), projectSlug, manifest.ProjectId),
		)

		if err := docker.RunContainerFromConfig(ctx, config); err != nil {
			return err
		}
		if err := dockerClient.ConnectContainerToNetwork(ctx, containerName, projectNetwork); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}

	// Manage docker apps in the manifest
	for _, app := range m.Apps.Docker {
		containerName := fmt.Sprintf("%s-%s-%s", app.Name, projectSlug, projectID)
		fmt.Printf("Deploying Docker image app: %s\n", app.Name)
		if err := dockerClient.RunContainer(ctx, containerName, app.Image); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
		if err := dockerClient.ConnectContainerToNetwork(ctx, containerName, projectNetwork); err != nil {
			return fmt.Errorf(errFmtAppRunFailed+": %w", app.Name, err)
		}
	}

	// Manage S3 services
	for _, svc := range m.Services.S3 {
		s3svc := services.NewS3Service(svc, s.cephEndpoint)
		deployed, err := s3svc.IsDeployed(ctx, projectID, *m)
		if err != nil {
			return err
		}
		if !deployed {
			if err := s3svc.Prepare(ctx, projectID, dockerClient, *m); err != nil {
				return fmt.Errorf("failed to prepare S3 service %q: %w", svc.Name, err)
			}
			if err := s3svc.Deploy(ctx, projectID, dockerClient, *m); err != nil {
				return fmt.Errorf("failed to deploy S3 service %q: %w", svc.Name, err)
			}
		}
	}

	return nil
}
