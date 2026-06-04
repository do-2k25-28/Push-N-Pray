package service

import (
	"context"
	"fmt"
	"pushnpray/internal"
)

func UpdateServices(ctx context.Context, dockerClient *internal.Client, projectID string, declaredServices []ServiceDefinition) error {
	existingServices, err := dockerClient.ListServiceContainers(ctx, projectID)
	if err != nil {
		return err
	}

	declaredByID := make(map[string]ServiceDefinition, len(declaredServices))
	for _, definition := range declaredServices {
		declaredByID[definition.ID] = definition
	}

	existingByID := make(map[string]internal.ServiceContainer, len(existingServices))
	for _, existingService := range existingServices {
		serviceID := existingService.Labels[labelID]
		existingByID[serviceID] = existingService
		if definition, declared := declaredByID[serviceID]; declared {
			if existingService.Labels[labelType] != definition.Type || existingService.Labels[labelVersion] != definition.Version {
				return fmt.Errorf("service %q cannot be modified after creation", serviceID)
			}
		}
	}

	for serviceID, existingService := range existingByID {
		if _, declared := declaredByID[serviceID]; declared {
			continue
		}
		if err := dockerClient.RemoveContainerAndVolume(ctx, existingService.Name, existingService.Labels[labelVolume]); err != nil {
			return fmt.Errorf("remove service %q: %w", serviceID, err)
		}
	}

	for _, definition := range declaredServices {
		if existingService, exists := existingByID[definition.ID]; exists {
			if !existingService.Running {
				if err := dockerClient.StartContainer(ctx, existingService.Name); err != nil {
					return fmt.Errorf("start service %q: %w", definition.ID, err)
				}
			}
			continue
		}

		if err := dockerClient.EnsureVolume(ctx, definition.volumeName(projectID), definition.labels(projectID)); err != nil {
			return fmt.Errorf("create volume for service %q: %w", definition.ID, err)
		}
		if err := dockerClient.RunContainer(ctx, definition.containerConfig(projectID)); err != nil {
			return fmt.Errorf("provision service %q: %w", definition.ID, err)
		}
	}

	return nil
}
