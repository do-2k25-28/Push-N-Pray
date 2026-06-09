package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	s3infra "pushnpray/infrastructure/s3"
	cephinfra "pushnpray/internal/ceph"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

const cephContainerName = "ceph"
const cephContainerEndpoint = "http://ceph:8080"

type S3Service struct {
	Manifest manifest.S3Service
	endpoint string
}

func NewS3Service(svc manifest.S3Service, endpoint string) *S3Service {
	return &S3Service{Manifest: svc, endpoint: endpoint}
}

func (s *S3Service) bucketName(projectID string) string {
	return fmt.Sprintf("%s-%s", s.Manifest.Name, projectID)
}

func (s *S3Service) userID(projectID string) string {
	return fmt.Sprintf("project-%s-%s", projectID, s.Manifest.Name)
}

func getS3ServiceData(projectID, name string) (*models.S3Service, error) {
	var svc models.S3Service
	if err := database.GetDB().Where("project = ? AND name = ?", projectID, name).First(&svc).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

func (s *S3Service) IsDeployed(ctx context.Context, m manifest.Manifest) (bool, error) {
	var v int64
	res := database.GetDB().Table("s3_services").Select("1").
		Where("project = ? AND name = ?", m.ProjectId, s.Manifest.Name).Limit(1).Find(&v)
	return res.RowsAffected > 0, res.Error
}

func (s *S3Service) Prepare(ctx context.Context, docker *dockerw.Client, m manifest.Manifest) error {
	accessKey := rand.Text()
	secretKey := rand.Text()

	if err := cephinfra.CreateUser(ctx, docker, s.userID(m.ProjectId), accessKey, secretKey); err != nil {
		return fmt.Errorf("create Ceph user: %w", err)
	}

	return database.GetDB().Create(&models.S3Service{
		Project:   m.ProjectId,
		Name:      s.Manifest.Name,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}).Error
}

func (s *S3Service) Deploy(ctx context.Context, docker *dockerw.Client, m manifest.Manifest, network dockerw.ContainerNetwork) error {
	data, err := getS3ServiceData(m.ProjectId, s.Manifest.Name)
	if err != nil {
		return err
	}

	cfg := s3infra.NewConfig(s.endpoint, data.AccessKey, data.SecretKey)
	if err := s3infra.NewClientFromConfig(cfg).EnsureBucketExists(ctx, s.bucketName(m.ProjectId)); err != nil {
		return err
	}

	return docker.ConnectContainerToNetwork(ctx, cephContainerName, network.Name)
}

func (s *S3Service) EnvToInject(m manifest.Manifest) (map[string]map[string]string, error) {
	data, err := getS3ServiceData(m.ProjectId, s.Manifest.Name)
	if err != nil {
		return nil, err
	}

	prefix := "S3_" + strings.ToUpper(s.Manifest.Name) + "_"
	env := map[string]string{
		prefix + "ENDPOINT":   cephContainerEndpoint,
		prefix + "ACCESS_KEY": data.AccessKey,
		prefix + "SECRET_KEY": data.SecretKey,
		prefix + "BUCKET":     s.bucketName(m.ProjectId),
	}

	labels := map[string]map[string]string{}
	for _, app := range m.Apps.Docker {
		labels[app.Name] = env
	}
	for _, app := range m.Apps.Dockerfile {
		labels[app.Name] = env
	}
	return labels, nil
}
