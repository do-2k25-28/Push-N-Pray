package services

import (
	"context"
	"fmt"
	"strings"

	"pushnpray/cmd/server/deployment/project"
	s3infra "pushnpray/infrastructure/s3"
	"pushnpray/internal"
	"pushnpray/internal/manifest"
)

const cephContainerName = "ceph"

type S3Service struct {
	Manifest manifest.S3Service
	client   *s3infra.Client
	config   s3infra.Config
}

func NewS3Service(svc manifest.S3Service, serverEndpoint string) *S3Service {
	cfg := s3infra.NewConfig(serverEndpoint, svc.AccessKey, svc.SecretKey)
	return &S3Service{
		Manifest: svc,
		client:   s3infra.NewClientFromConfig(cfg),
		config:   cfg,
	}
}

func (s *S3Service) bucketName(projectID string) string {
	return fmt.Sprintf("%s-%s", s.Manifest.Name, projectID)
}

func (s *S3Service) IsDeployed(ctx context.Context, projectID string, m manifest.Manifest) (bool, error) {
	return s.client.BucketExists(ctx, s.bucketName(projectID))
}

func (s *S3Service) Prepare(_ context.Context, _ string, _ *internal.Client, _ manifest.Manifest) error {
	return nil
}

// Deploy creates the S3 bucket and connects the Ceph container to the project
// network so that app containers can reach it via http://ceph:8080.
func (s *S3Service) Deploy(ctx context.Context, projectID string, docker *internal.Client, m manifest.Manifest) error {
	if err := s.client.EnsureBucketExists(ctx, s.bucketName(projectID)); err != nil {
		return err
	}
	networkName := project.NetworkName(projectID)
	if err := docker.EnsureNetworkExists(ctx, networkName); err != nil {
		return err
	}
	return docker.ConnectContainerToNetwork(ctx, cephContainerName, networkName)
}

func (s *S3Service) EnvToInject(projectID string, m manifest.Manifest) (map[string]map[string]string, error) {
	labels := map[string]map[string]string{}
	bucket := s.bucketName(projectID)
	prefix := "S3_" + strings.ToUpper(s.Manifest.Name) + "_"

	env := map[string]string{
		prefix + "ENDPOINT":   s.config.ContainerEndpoint,
		prefix + "ACCESS_KEY": s.config.AccessKey,
		prefix + "SECRET_KEY": s.config.SecretKey,
		prefix + "BUCKET":     bucket,
	}

	for _, app := range m.Apps.Docker {
		labels[app.Name] = env
	}
	for _, app := range m.Apps.Dockerfile {
		labels[app.Name] = env
	}

	return labels, nil
}
