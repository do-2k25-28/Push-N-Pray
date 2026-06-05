package services

import (
	"context"
	"pushnpray/internal/dockerw"
	"pushnpray/internal/manifest"
)

type ManagedService interface {
	// Check if service is already deployed
	// If this returns true, service will be deployed
	IsDeployed(ctx context.Context, manifest manifest.Manifest) (bool, error)
	// Create volumes, idk can be anything
	// It's called once the very first time the service is deployed
	Prepare(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error
	// Deploy the service
	Deploy(ctx context.Context, client *dockerw.Client, manifest manifest.Manifest) error
	// Environment variables to inject
	// This is called everytime the project is deployed
	EnvToInject(manifest manifest.Manifest) (map[string]map[string]string, error)
}
