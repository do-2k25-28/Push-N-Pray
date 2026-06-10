package project

import (
	"fmt"
	"log"
	"os"
)

func TraefikNet() string {
	net := os.Getenv("TRAEFIK_NET")

	if net == "" {
		log.Fatal("Missing TRAEFIK_NET environment variable.")
	}

	return net
}

func domainSuffix() string {
	if s := os.Getenv("DOMAIN_SUFFIX"); s != "" {
		return s
	}
	return "localhost"
}

func TraefikLabels(containerName, appName, projectSlug, projectID string) map[string]string {
	suffix := domainSuffix()
	domain := fmt.Sprintf("%s-%s-%s.%s", appName, projectSlug, projectID, suffix)

	labels := map[string]string{
		"traefik.enable":         "true",
		"traefik.docker.network": TraefikNet(),
		fmt.Sprintf("traefik.http.routers.%s.rule", containerName): fmt.Sprintf("Host(`%s`)", domain),
	}

	if suffix == "localhost" {
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName)] = "web"
	} else {
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName)] = "websecure"
		labels[fmt.Sprintf("traefik.http.routers.%s.tls", containerName)] = "true"
		labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", containerName)] = "le"
	}

	return labels
}
