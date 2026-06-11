package updates

import (
	"context"
	"fmt"
	"log"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type UpdateStrategy interface {
	UpdateContainer(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string) error
}

func RunApplicationUpdate(ctx context.Context, docker *dockerw.Client, newContainer dockerw.ContainerConfig, pattern string, strategyName manifest.UpdateStrategy) error {
	var strategy UpdateStrategy = nil

	if strategyName == manifest.UpdateStrategyRecreate {
		strategy = RecreateUpdateStrategy{}
	}

	if strategyName == manifest.UpdateStrategyRolling {
		strategy = RollingUpdateStrategy{}
	}

	if strategyName == manifest.UpdateStrategyBlueGreen {
		strategy = BlueGreenUpdateStrategy{}
	}

	if strategyName == manifest.UpdateStrategyCanary {
		strategy = CanaryUpdateStrategy{}
	}

	if strategy == nil {
		return fmt.Errorf("unknown app update strategy %s", strategyName)
	}

	log.Printf("Using update strategy %s\n", strategyName)

	return strategy.UpdateContainer(ctx, docker, newContainer, pattern)
}
