package deployment

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"path/filepath"

	dockerSdk "github.com/docker/go-sdk/client"
	sdkcontainer "github.com/docker/go-sdk/container"
	sdkimage "github.com/docker/go-sdk/image"
	sdknetwork "github.com/docker/go-sdk/network"
	"github.com/moby/moby/api/types/container"
	mobytypesnetwork "github.com/moby/moby/api/types/network"
	dockerclient "github.com/moby/moby/client"
)

const (
	traefikName  = "pushnpray-traefik"
	traefikNet   = "traefik"
	traefikImage = "traefik:v3.7"
)

func EnsureTraefik(ctx context.Context) error {
	cli, err := dockerSdk.New(ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", errFmtTraefikClient, err)
	}

	if _, err = sdknetwork.New(ctx, sdknetwork.WithName(traefikNet), sdknetwork.WithClient(cli)); err != nil {
		return fmt.Errorf("%s: %w", errFmtTraefikNetwork, err)
	}

	acmePath := filepath.Join("/var/lib/pushnpray", "traefik", "acme.json")
	_ = os.MkdirAll(filepath.Dir(acmePath), 0700)
	_ = os.WriteFile(acmePath, []byte("{}"), 0600)
	_ = os.Chmod(acmePath, 0600)

	return ensureTraefikContainer(ctx, cli, acmePath)
}

func ensureTraefikContainer(ctx context.Context, cli dockerSdk.SDKClient, acmePath string) error {
	if _, err := cli.ContainerStart(ctx, traefikName, dockerclient.ContainerStartOptions{}); err == nil {
		return nil
	}

	if err := sdkimage.Pull(ctx, traefikImage, sdkimage.WithPullClient(cli)); err != nil {
		return fmt.Errorf("%s: %w", errFmtTraefikPull, err)
	}

	email := os.Getenv("TRAEFIK_ACME_EMAIL")
	if email == "" {
		email = "admin@pushnpray.polydo.dev"
	}

	_, err := sdkcontainer.Run(
		ctx,
		sdkcontainer.WithClient(cli),
		sdkcontainer.WithName(traefikName),
		sdkcontainer.WithImage(traefikImage),
		sdkcontainer.WithNetworkName(nil, traefikNet),
		sdkcontainer.WithExposedPorts("80/tcp", "443/tcp"),
		sdkcontainer.WithCmd(
			"--providers.docker=true",
			"--providers.docker.network="+traefikNet,
			"--providers.docker.exposedbydefault=false",
			"--entrypoints.web.address=:80",
			"--entrypoints.websecure.address=:443",
			"--certificatesresolvers.le.acme.email="+email,
			"--certificatesresolvers.le.acme.storage=/acme.json",
			"--certificatesresolvers.le.acme.httpchallenge=true",
			"--certificatesresolvers.le.acme.httpchallenge.entrypoint=web",
		),
		sdkcontainer.WithHostConfigModifier(func(h *container.HostConfig) {
			h.RestartPolicy = container.RestartPolicy{Name: "unless-stopped"}
			h.PortBindings = mobytypesnetwork.PortMap{
				mobytypesnetwork.MustParsePort("80/tcp"):  []mobytypesnetwork.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "80"}},
				mobytypesnetwork.MustParsePort("443/tcp"): []mobytypesnetwork.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "443"}},
			}
			h.Binds = []string{"/var/run/docker.sock:/var/run/docker.sock:ro", acmePath + ":/acme.json"}
		}),
	)

	if err != nil {
		return fmt.Errorf("%s: %w", errFmtTraefikRun, err)
	}

	return nil
}
