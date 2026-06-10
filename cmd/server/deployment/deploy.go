package deployment

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/apps"
	"pushnpray/cmd/server/deployment/project"
	"pushnpray/cmd/server/deployment/services"
	"pushnpray/internal/ceph"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
	"pushnpray/internal/utils"
)

func DeployProject(projectSlug string, projectManifest manifest.Manifest, workspaceDir, deploymentID string) error {
	ctx := context.WithValue(context.Background(), apps.WorkingDirectoryContextKey, workspaceDir)
	docker, err := dockerw.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	if projectManifest.AppUpdateStrategy == "" {
		projectManifest.AppUpdateStrategy = manifest.AppUpdateStrategyRecreate
	}
	if projectManifest.AppUpdateStrategy != manifest.AppUpdateStrategyRecreate && projectManifest.AppUpdateStrategy != manifest.AppUpdateStrategyBlueGreen {
		return fmt.Errorf("unknown app update strategy %q", projectManifest.AppUpdateStrategy)
	}

	// Create project Docker network

	if err := docker.CreateNetworkIfNotExist(ctx, project.NetworkName(projectManifest.ProjectId)); err != nil {
		return fmt.Errorf("failed to create project network: %w", err)
	}

	network := dockerw.ContainerNetwork{
		Name:    project.NetworkName(projectManifest.ProjectId),
		Aliases: []string{},
	}
	traefikNet := dockerw.ContainerNetwork{
		Name:    project.TraefikNet(),
		Aliases: []string{},
	}

	// Handle managed services creation and deletion

	_services := make([]services.ManagedService, 0, projectManifest.GetServiceCount())
	envsFromServices := []map[string]map[string]string{}

	for _, service := range projectManifest.Services.Postgres {
		pg := services.PostgresService{Manifest: service}
		_services = append(_services, &pg)
	}
	for _, service := range projectManifest.Services.S3 {
		_services = append(_services, services.NewS3Service(service, ceph.GetCephEndpoint()))
	}

	for _, service := range projectManifest.Services.Redis {
		redis := services.RedisService{Manifest: service}
		_services = append(_services, &redis)
	}

	for _, service := range _services {
		deployed, err := service.IsDeployed(ctx, projectManifest)
		if err != nil {
			return err
		}

		if !deployed {
			if err := service.Prepare(ctx, docker, projectManifest); err != nil {
				return fmt.Errorf("filed to prepare deployment of service")
			}

			if err := service.Deploy(ctx, docker, projectManifest, network); err != nil {
				return fmt.Errorf("failed to deploy service")
			}
		}

		env, err := service.EnvToInject(projectManifest)
		if err != nil {
			return err
		}
		envsFromServices = append(envsFromServices, env)
	}

	appToEnv := utils.MergeMaps(envsFromServices)

	// Deploy or update application containers
	_apps := make([]apps.DeployableApp, 0, projectManifest.GetApplicationCount())

	for _, app := range projectManifest.Apps.Docker {
		_apps = append(_apps, apps.NewDockerApp(app))
	}

	for _, app := range projectManifest.Apps.Dockerfile {
		_apps = append(_apps, apps.NewDockerFileApp(app, workspaceDir))
	}

	for _, app := range projectManifest.Apps.StaticWeb {
		_apps = append(_apps, apps.NewStaticWebApp(app))
	}

	for _, app := range _apps {
		if err := app.Prepare(ctx, docker, projectManifest); err != nil {
			return err
		}

		config := app.ContainerConfig(ctx, projectManifest)

		config.Env = utils.MergeMap(
			config.Env,
			appToEnv[app.AppName()], // Override user defined vars if they overlap
		)

		baseName := "app-" + projectManifest.ProjectId + "-" + app.AppName()
		config.Name = baseName
		
		if projectManifest.AppUpdateStrategy == manifest.AppUpdateStrategyBlueGreen {
			config.Name = "app-" + deploymentID + "-" + projectManifest.ProjectId + "-" + app.AppName()
		}
		
		config.Networks = []dockerw.ContainerNetwork{network, traefikNet}
		config.Labels = utils.MergeMap(
			config.Labels,
			project.TraefikLabels(config.Name, app.AppName(), projectSlug, projectManifest.ProjectId),
		)

		if err := deployAppContainer(ctx, docker, config, baseName, projectManifest.AppUpdateStrategy); err != nil {
			return err
		}
	}

	return nil
}

func deployAppContainer(ctx context.Context, docker *dockerw.Client, config dockerw.ContainerConfig, baseName string, strategy manifest.AppUpdateStrategy) error {
	switch strategy {
	case manifest.AppUpdateStrategyRecreate:
		if err := docker.StopContainersByPattern(ctx, baseName); err != nil {
			return err
		}

		if err := docker.RemoveContainersByPattern(ctx, baseName); err != nil {
			return err
		}

		return docker.RunContainerFromConfig(ctx, config)
	case manifest.AppUpdateStrategyBlueGreen:
		if err := docker.RunContainerFromConfig(ctx, config); err != nil {
			return err
		}

		if err := docker.StopContainersByPattern(ctx, baseName); err != nil {
			return err
		}

		return docker.RemoveContainersByPattern(ctx, baseName)
	default:
		return fmt.Errorf("unknown app update strategy %q", strategy)
	}
}
