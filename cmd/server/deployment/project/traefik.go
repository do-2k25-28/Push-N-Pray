package project

import (
	"fmt"
	"log"
	"os"
)

func traefikNet() string {
	net := os.Getenv("TRAEFIK_NET")

	if net == "" {
		log.Fatal("Missing TRAEFIK_NET environment variable.")
	}

	return net
}

func TraefikLabels(containerName, appName, projectSlug, projectID string) map[string]string {
	domain := fmt.Sprintf("%s-%s-%s.pushnpray.polydo.dev", appName, projectSlug, projectID)
	return map[string]string{
		"traefik.enable":         "true",
		"traefik.docker.network": traefikNet(),
		fmt.Sprintf("traefik.http.routers.%s.rule", containerName):             fmt.Sprintf("Host(`%s`)", domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName):      "websecure",
		fmt.Sprintf("traefik.http.routers.%s.tls", containerName):              "true",
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", containerName): "le",
	}
}
