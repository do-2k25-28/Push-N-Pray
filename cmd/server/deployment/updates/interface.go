package updates

import (
	"context"
	"fmt"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type UpdateStrategy interface {
	UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error
}

func RunApplicationUpdate(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string, strategyName manifest.AppUpdateStrategy) error {
	var strategy UpdateStrategy = nil

	if strategyName == manifest.AppUpdateStrategyRecreate {
		strategy = RecreateUpdateStrategy{}
	}

	if strategyName == manifest.AppUpdateStrategyBlueGreen {
		strategy = BlueGreenUpdateStrategy{}
	}

	if strategy == nil {
		return fmt.Errorf("unknown app update strategy %q", strategy)
	}

	return strategy.UpdateContainer(ctx, docker, newContainer, pattern)
}
