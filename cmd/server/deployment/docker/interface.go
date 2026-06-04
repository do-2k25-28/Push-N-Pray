package docker

import (
	"context"
	"pushnpray/cmd/server/deployment/container"
)

type Client interface {
	container.Runner
	BuildImage(ctx context.Context, tag, dockerfilePath, contextDir string) error
}

type DeployableApp interface {
	RunContainer(ctx context.Context, client Client) error
}
