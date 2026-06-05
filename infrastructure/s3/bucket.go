package s3

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type Client struct {
	s3 *s3.Client
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
		BaseEndpoint:       aws.String(endpoint),
		Credentials:        creds,
		Region:             "us-east-1",
		UsePathStyle:       true,
	})

	return &Client{s3: client}, nil
}

// EnsureBucket creates the bucket if it does not already exist.
func (c *Client) EnsureBucket(ctx context.Context, bucketName string) error {
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
