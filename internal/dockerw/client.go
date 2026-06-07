package dockerw

import (
	"context"
	"fmt"

	"github.com/docker/go-sdk/client"
)

type Client struct {
	client.SDKClient
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := client.New(ctx)

	if err != nil {
		return nil, fmt.Errorf("dockerwrapper: create client: %w", err)
	}

	return &Client{client}, nil
}
