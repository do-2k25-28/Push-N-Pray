package deployment

import (
	"context"
	"fmt"
	"pushnpray/cmd/server/deployment/project"
	"pushnpray/cmd/server/deployment/services"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

func DeployProject(projectSlug string, projectID string, manifest manifest.Manifest, workspaceDir string) error {
	ctx := context.Background()
	docker, err := dockerw.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create project Docker network

	if err := docker.CreateNetworkIfNotExist(ctx, project.NetworkName(projectID)); err != nil {
		return fmt.Errorf("failed to create project network: %w", err)
	}

	// Handle managed services creation and deletion

	_services := make([]services.ManagedService, 0, manifest.GetServiceCount())
	envs := []map[string]map[string]string{}

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

			if err := service.Deploy(ctx, docker, manifest); err != nil {
				return fmt.Errorf("failed to deploy service")
			}
		}

		env, err := service.EnvToInject(manifest)
		if err != nil {
			return err
		}
		envs = append(envs, env)
	}

	fmt.Println(envs)

	// Deploy or update application containers

	return nil
}
