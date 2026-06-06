package s3

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"pushnpray/internal/manifest"
)

// Service provisions S3 resources for a specific project.
type Service interface {
	ProvisionBuckets(ctx context.Context, services []manifest.S3Service) error
}

type Client struct {
	s3 *s3.Client
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

func NewClient() (*Client, error) {
	endpoint := os.Getenv("CEPH_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	accessKey := os.Getenv("CEPH_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "demo-access-key"
	}
	secretKey := os.Getenv("CEPH_SECRET_KEY")
	if secretKey == "" {
		secretKey = "demo-secret-key"
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(endpoint),
		Credentials:  creds,
		Region:       "us-east-1",
		UsePathStyle: true,
	})

	return &Client{s3: client}, nil
}

func (c *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := c.s3.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
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
