package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
	"slices"
	"strings"
)

const postgresUser = "postgres"

type PostgresService struct {
	Manifest manifest.PostgresService
}

func (s *PostgresService) volumeName(projectId string) string {
	return "postgres-" + s.Manifest.Name + "-" + projectId
}

func getServiceDataFromDatabase(projectId string, name string) (*models.PostgresService, error) {
	var service models.PostgresService
	if err := database.GetDB().Where("project = ? AND name = ?", projectId, name).First(&service).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (s *PostgresService) IsDeployed(ctx context.Context, manifest manifest.Manifest) (bool, error) {
	var v int64
	res := database.GetDB().Table("postgres_services").Select("1").Where("project = ? AND name = ?", manifest.ProjectId, s.Manifest.Name).Limit(1).Find(&v)

	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

// Only thing required by the postgres service is a persistent docker volume and a set of credentials
func (s *PostgresService) Prepare(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error {
	// Create volume for postgres
	err := client.CreateVolumeIfNotExist(ctx, s.volumeName(manifest.ProjectId), map[string]string{})
	if err != nil {
		return err
	}

	// Create credentials and store them into postgres

	password := rand.Text()
	return database.GetDB().Create(models.PostgresService{
		Project:  manifest.ProjectId,
		Name:     s.Manifest.Name,
		Password: password,
	}).Error
}

func (s *PostgresService) Deploy(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error {
	fmt.Printf("%+v\n", s)
	fmt.Printf("%+v\n", manifest)

	data, err := getServiceDataFromDatabase(manifest.ProjectId, s.Manifest.Name)
	if err != nil {
		return err
	}

	container := dockerw.ContainerConfig{
		Image: "docker.io/library/postgres:18.4-alpine3.23",
		Name:  "postgres-" + manifest.ProjectId + "-" + s.Manifest.Name,
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": data.Password,
		},
		VolumeBinds: []string{s.volumeName(manifest.ProjectId) + ":/var/lib/postgresql"},
	}

	return client.RunContainerFromConfig(ctx, container)
}

func (s *PostgresService) EnvToInject(manifest manifest.Manifest) (map[string]map[string]string, error) {
	fmt.Printf("%+v\n", s)
	fmt.Printf("%+v\n", manifest)

	data, err := getServiceDataFromDatabase(manifest.ProjectId, s.Manifest.Name)
	if err != nil {
		return nil, err
	}

	labels := map[string]map[string]string{}

	for _, app := range manifest.GetApps() {
		if slices.Contains(s.Manifest.UsedBy, app.Name) {
			prefix := "POSTGRES_" + strings.ToUpper(app.Name) + "_"

			labels[app.Name] = map[string]string{
				prefix + "USER":     postgresUser,
				prefix + "PASSWORD": data.Password,
				prefix + "HOST":     "postgres-" + app.Name, // TODO: check if host is resolved by Docker internal DNS
			}
		}
	}

	return labels, nil
}
