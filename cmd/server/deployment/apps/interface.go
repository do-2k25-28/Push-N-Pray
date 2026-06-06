package apps

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type DeployableApp interface {
	AppName() string
	// Prepare function is run before getting the contaienr config
	// Can be anything. For exemple building container images
	Prepare(ctx context.Context, docker *dockerw.Client, manifest manifest.Manifest) error
	// Basic container config
	// Container name and network is managed by the deploy function not this
	// Env can be populated but the deploy engine will add managed services
	// environment variables
	ContainerConfig(ctx context.Context, manifest manifest.Manifest) dockerw.ContainerConfig
}
