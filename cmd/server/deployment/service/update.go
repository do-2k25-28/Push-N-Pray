package service

import (
	"context"
	"fmt"
	"pushnpray/internal"
)

func UpdateServices(ctx context.Context, dockerClient *internal.Client, projectId string, declaredServices []ServiceDefinition) error {
	existingServices, err := dockerClient.ListServiceContainers(ctx, projectId)
	if err != nil {
		return err
	}

	declaredById := make(map[string]ServiceDefinition, len(declaredServices))
	for _, definition := range declaredServices {
		declaredById[definition.Id] = definition
	}

	existingById := make(map[string]internal.ServiceContainer, len(existingServices))
	for _, existingService := range existingServices {
		serviceId := existingService.Labels[labelId]
		existingById[serviceId] = existingService
		if definition, declared := declaredById[serviceId]; declared {
			if existingService.Labels[labelType] != definition.Type || existingService.Labels[labelVersion] != definition.Version {
				return fmt.Errorf("service %q cannot be modified after creation", serviceId)
			}
		}
	}

	for serviceId, existingService := range existingById {
		if _, declared := declaredById[serviceId]; declared {
			continue
		}
		if err := dockerClient.RemoveContainerAndVolume(ctx, existingService.Name, existingService.Labels[labelVolume]); err != nil {
			return fmt.Errorf("remove service %q: %w", serviceId, err)
		}
	}

	for _, definition := range declaredServices {
		if existingService, exists := existingById[definition.Id]; exists {
			if !existingService.Running {
				if err := dockerClient.StartContainer(ctx, existingService.Name); err != nil {
					return fmt.Errorf("start service %q: %w", definition.Id, err)
				}
			}
			continue
		}

		if err := dockerClient.EnsureVolume(ctx, definition.volumeName(projectId), definition.labels(projectId)); err != nil {
			return fmt.Errorf("create volume for service %q: %w", definition.Id, err)
		}
		if err := dockerClient.RunContainer(ctx, definition.containerConfig(projectId)); err != nil {
			return fmt.Errorf("provision service %q: %w", definition.Id, err)
		}
	}

	return nil
}
