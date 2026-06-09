package deployment

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/apps"
	"pushnpray/cmd/server/deployment/project"
	"pushnpray/cmd/server/deployment/services"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
	"pushnpray/internal/utils"
)

func DeployProject(projectSlug string, manifest manifest.Manifest, workspaceDir string) error {
	ctx := context.WithValue(context.Background(), apps.WorkingDirectoryContextKey, workspaceDir)
	docker, err := dockerw.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create project Docker network

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

	for _, service := range manifest.Services.Redis {
		redis := services.RedisService{Manifest: service}
		_services = append(_services, &redis)
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

	for _, app := range manifest.Apps.StaticWeb {
		_apps = append(_apps, apps.NewStaticWebApp(app))
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
	}

	return nil
}
