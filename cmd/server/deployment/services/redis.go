package services

import (
	"context"
	"crypto/rand"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
	"slices"
	"strings"
)

type RedisService struct {
	Manifest manifest.RedisService
}

func (s RedisService) containerName(projectId string) string {
	return "redis-" + projectId + "-" + s.Manifest.Name
}

func getRedisDataFromDatabase(projectId string, name string) (*models.RedisService, error) {
	var service models.RedisService
	if err := database.GetDB().Where("project = ? AND name = ?", projectId, name).First(&service).Error; err != nil {
		return nil, err
	}

	return &service, nil
}

func (s RedisService) IsDeployed(ctx context.Context, manifest manifest.Manifest) (bool, error) {
	// TODO!: use container presence as a source of truth, not the database

	var v int64
	res := database.GetDB().Table("redis_services").Select("1").Where("project = ? AND name = ?", manifest.ProjectId, s.Manifest.Name).Limit(1).Find(&v)

	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

// Prepare generates and stores credentials for the redis service
func (s RedisService) Prepare(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error {
	password := rand.Text()
	return database.GetDB().Create(models.RedisService{
		Project:  manifest.ProjectId,
		Name:     s.Manifest.Name,
		Password: password,
	}).Error
}

func (s RedisService) Deploy(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest, network dockerw.ContainerNetwork) error {
	data, err := getRedisDataFromDatabase(manifest.ProjectId, s.Manifest.Name)
	if err != nil {
		return err
	}

	container := dockerw.ContainerConfig{
		Image: "docker.io/bitnami/redis@sha256:6e7a020f1f6504698a7272c58783bdc2c23588c49febbae5aca1bb8dfa10af25",
		Name:  s.containerName(manifest.ProjectId),
		Env: map[string]string{
			"REDIS_DATABASE":            "pushnpray",
			"REDIS_PASSWORD":            data.Password,
			"REDIS_AOF_ENABLED":         "no",
			"REDIS_RDB_POLICY_DISABLED": "yes",
		},
		Networks: []dockerw.ContainerNetwork{network},
	}

	return client.RunContainerFromConfig(ctx, container)
}

func (s RedisService) EnvToInject(manifest manifest.Manifest) (map[string]map[string]string, error) {
	data, err := getRedisDataFromDatabase(manifest.ProjectId, s.Manifest.Name)
	if err != nil {
		return nil, err
	}

	labels := map[string]map[string]string{}

	for _, app := range manifest.GetApps() {
		if slices.Contains(s.Manifest.UsedBy, app.Name) {
			prefix := "REDIS_" + strings.ToUpper(s.Manifest.Name) + "_"

			labels[app.Name] = map[string]string{
				prefix + "USER":     "default",
				prefix + "PASSWORD": data.Password,
				prefix + "HOST":     s.containerName(manifest.ProjectId),
				prefix + "PORT":     "6379",
			}
		}
	}

	return labels, nil
}
