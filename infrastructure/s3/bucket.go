package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"pushnpray/internal/manifest"
)

type Config struct {
	Endpoint          string
	ContainerEndpoint string
	AccessKey         string
	SecretKey         string
}

func NewConfig(serverEndpoint, accessKey, secretKey string) Config {
	return Config{
		Endpoint:          serverEndpoint,
		ContainerEndpoint: "http://ceph:8080",
		AccessKey:         accessKey,
		SecretKey:         secretKey,
	}
}

type Service interface {
	ProvisionBuckets(ctx context.Context, services []manifest.S3Service) error
}

type Client struct {
	s3 *s3.Client
}

// NewClientFromConfig creates an S3 client from the given config.
func NewClientFromConfig(cfg Config) *Client {
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")
	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(cfg.Endpoint),
		Credentials:  creds,
		Region:       "us-east-1",
		UsePathStyle: true,
	})
	return &Client{s3: client}
}

// ProjectService scopes S3 operations to a single project.
type ProjectService struct {
	client      *Client
	projectSlug string
	projectID   string
}

func NewProjectService(client *Client, projectSlug, projectID string) *ProjectService {
	return &ProjectService{client: client, projectSlug: projectSlug, projectID: projectID}
}

func (s *ProjectService) bucketName(serviceName string) string {
	return fmt.Sprintf("%s-%s-%s", serviceName, s.projectSlug, s.projectID)
}

func (s *ProjectService) ProvisionBuckets(ctx context.Context, services []manifest.S3Service) error {
	for _, svc := range services {
		name := s.bucketName(svc.Name)
		if err := s.client.EnsureBucketExists(ctx, name); err != nil {
			return fmt.Errorf("failed to provision S3 bucket for service %q: %w", svc.Name, err)
		}
		fmt.Printf("S3 bucket %q ready\n", name)
	}
	return nil
}

func (c *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := c.s3.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

// BucketExists reports whether the given bucket exists and is accessible.
func (c *Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return true, nil
	}
	var notFound *s3types.NotFound
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, fmt.Errorf("check bucket %q: %w", bucketName, err)
}

// EnsureBucketExists creates the bucket if it does not already exist.
func (c *Client) EnsureBucketExists(ctx context.Context, bucketName string) error {
	_, err := c.s3.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(smithy.APIError); ok {
		switch apiErr.ErrorCode() {
		case "BucketAlreadyExists", "BucketAlreadyOwnedByYou":
			return nil
		}
	}

	return fmt.Errorf("create bucket %q: %w", bucketName, err)
}
