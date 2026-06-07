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

type Client struct {
	s3 *s3.Client
}

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

func (c *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := c.s3.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}
