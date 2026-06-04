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
	labelProjectID = "pushnpray.project-id"
	labelID        = "pushnpray.service-id"
	labelType      = "pushnpray.service-type"
	labelVersion   = "pushnpray.service-version"
	labelVolume    = "pushnpray.service-volume"
)

var validServiceID = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

type ServiceDefinition struct {
	ID      string
	Type    string
	Version string
}

func ServiceDefinitionsFromManifest(projectManifest *manifest.Manifest) ([]ServiceDefinition, error) {
	definitions := make([]ServiceDefinition, 0, len(projectManifest.Services.Postgres)+len(projectManifest.Services.Redis)+len(projectManifest.Services.S3))
	for _, service := range projectManifest.Services.Postgres {
		definitions = append(definitions, ServiceDefinition{ID: service.ID, Type: "postgres", Version: service.Version})
	}
	for _, service := range projectManifest.Services.Redis {
		definitions = append(definitions, ServiceDefinition{ID: service.ID, Type: "redis", Version: service.Version})
	}
	for _, service := range projectManifest.Services.S3 {
		definitions = append(definitions, ServiceDefinition{ID: service.ID, Type: "s3", Version: "latest"})
	}

	usedIDs := make(map[string]bool, len(definitions))
	for _, definition := range definitions {
		if definition.ID == "" {
			return nil, fmt.Errorf("%s service id is required", definition.Type)
		}
		if !validServiceID.MatchString(definition.ID) {
			return nil, fmt.Errorf("service id %q must contain only lowercase letters, numbers, and hyphens", definition.ID)
		}
		if definition.Version == "" {
			return nil, fmt.Errorf("%s service %q version is required", definition.Type, definition.ID)
		}
		if usedIDs[definition.ID] {
			return nil, fmt.Errorf("service id %q is duplicated", definition.ID)
		}
		usedIDs[definition.ID] = true
	}
	return definitions, nil
}

func (definition ServiceDefinition) containerName(projectID string) string {
	return fmt.Sprintf("service-%s-%s-%s", definition.Type, definition.ID, projectID)
}

func (definition ServiceDefinition) volumeName(projectID string) string {
	return definition.containerName(projectID) + "-data"
}

func (definition ServiceDefinition) labels(projectID string) map[string]string {
	return map[string]string{
		labelService:   "true",
		labelProjectID: projectID,
		labelID:        definition.ID,
		labelType:      definition.Type,
		labelVersion:   definition.Version,
		labelVolume:    definition.volumeName(projectID),
	}
}

func (definition ServiceDefinition) containerConfig(projectID string) internal.ContainerConfig {
	config := internal.ContainerConfig{
		Name:     definition.containerName(projectID),
		Networks: []internal.ContainerNetwork{{Name: container.ProjectNetworkName(projectID), Aliases: []string{definition.ID}}},
		Labels:   definition.labels(projectID),
	}
	volumeName := definition.volumeName(projectID)

	switch definition.Type {
	case "postgres":
		config.Image = "postgres:" + definition.Version
		config.Environment = map[string]string{
			"POSTGRES_DB":       "pushnpray",
			"POSTGRES_USER":     "pushnpray",
			"POSTGRES_PASSWORD": serviceSecret(projectID, definition.ID),
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
			"MINIO_ROOT_PASSWORD": serviceSecret(projectID, definition.ID),
		}
		config.VolumeBinds = []string{volumeName + ":/data"}
		config.Command = []string{"server", "/data"}
	}
	return config
}

func serviceSecret(projectID, serviceID string) string {
	sum := sha256.Sum256([]byte(projectID + ":" + serviceID))
	return hex.EncodeToString(sum[:16])
}
