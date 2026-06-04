package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"pushnpray/cmd/server/deployment/container"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
	"regexp"
)

const (
	labelService   = "pushnpray.service"
	labelProjectId = "pushnpray.project-id"
	labelId        = "pushnpray.service-id"
	labelType      = "pushnpray.service-type"
	labelVersion   = "pushnpray.service-version"
	labelVolume    = "pushnpray.service-volume"
)

var validServiceId = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

type ServiceDefinition struct {
	Id      string
	Type    string
	Version string
}

func ServiceDefinitionsFromManifest(projectManifest *manifest.Manifest) ([]ServiceDefinition, error) {
	definitions := make([]ServiceDefinition, 0, len(projectManifest.Services.Postgres)+len(projectManifest.Services.Redis)+len(projectManifest.Services.S3))
	for _, service := range projectManifest.Services.Postgres {
		definitions = append(definitions, ServiceDefinition{Id: service.Id, Type: "postgres", Version: service.Version})
	}
	for _, service := range projectManifest.Services.Redis {
		definitions = append(definitions, ServiceDefinition{Id: service.Id, Type: "redis", Version: service.Version})
	}
	for _, service := range projectManifest.Services.S3 {
		definitions = append(definitions, ServiceDefinition{Id: service.Id, Type: "s3", Version: "latest"})
	}

	usedIds := make(map[string]bool, len(definitions))
	for _, definition := range definitions {
		if definition.Id == "" {
			return nil, fmt.Errorf("%s service id is required", definition.Type)
		}
		if !validServiceId.MatchString(definition.Id) {
			return nil, fmt.Errorf("service id %q must contain only lowercase letters, numbers, and hyphens", definition.Id)
		}
		if definition.Version == "" {
			return nil, fmt.Errorf("%s service %q version is required", definition.Type, definition.Id)
		}
		if usedIds[definition.Id] {
			return nil, fmt.Errorf("service id %q is duplicated", definition.Id)
		}
		usedIds[definition.Id] = true
	}
	return definitions, nil
}

func (definition ServiceDefinition) containerName(projectId string) string {
	return fmt.Sprintf("service-%s-%s-%s", definition.Type, definition.Id, projectId)
}

func (definition ServiceDefinition) volumeName(projectId string) string {
	return definition.containerName(projectId) + "-data"
}

func (definition ServiceDefinition) labels(projectId string) map[string]string {
	return map[string]string{
		labelService:   "true",
		labelProjectId: projectId,
		labelId:        definition.Id,
		labelType:      definition.Type,
		labelVersion:   definition.Version,
		labelVolume:    definition.volumeName(projectId),
	}
}

func (definition ServiceDefinition) containerConfig(projectId string) internal.ContainerConfig {
	config := internal.ContainerConfig{
		Name:     definition.containerName(projectId),
		Networks: []internal.ContainerNetwork{{Name: container.ProjectNetworkName(projectId), Aliases: []string{definition.Id}}},
		Labels:   definition.labels(projectId),
	}
	volumeName := definition.volumeName(projectId)

	switch definition.Type {
	case "postgres":
		config.Image = "postgres:" + definition.Version
		config.Environment = map[string]string{
			"POSTGRES_DB":       "pushnpray",
			"POSTGRES_USER":     "pushnpray",
			"POSTGRES_PASSWORD": serviceSecret(projectId, definition.Id),
			"PGDATA":            "/var/lib/postgresql/data",
		}
		config.VolumeBinds = []string{volumeName + ":/var/lib/postgresql"}
	case "redis":
		config.Image = "redis:" + definition.Version
		config.VolumeBinds = []string{volumeName + ":/data"}
		config.Command = []string{"redis-server", "--appendonly", "yes"}
	case "s3":
		config.Image = "minio/minio:latest"
		config.Environment = map[string]string{
			"MINIO_ROOT_USER":     "pushnpray",
			"MINIO_ROOT_PASSWORD": serviceSecret(projectId, definition.Id),
		}
		config.VolumeBinds = []string{volumeName + ":/data"}
		config.Command = []string{"server", "/data"}
	}
	return config
}

func serviceSecret(projectId, serviceId string) string {
	sum := sha256.Sum256([]byte(projectId + ":" + serviceId))
	return hex.EncodeToString(sum[:16])
}
