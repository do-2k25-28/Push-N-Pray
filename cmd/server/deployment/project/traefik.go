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

func TraefikLabels(containerName, appName, projectSlug, projectID string) map[string]string {
	domain := fmt.Sprintf("%s-%s-%s.pushnpray.polydo.dev", appName, projectSlug, projectID)
	return map[string]string{
		"traefik.enable":         "true",
		"traefik.docker.network": TraefikNet(),
		fmt.Sprintf("traefik.http.routers.%s.rule", containerName):             fmt.Sprintf("Host(`%s`)", domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints", containerName):      "websecure",
		fmt.Sprintf("traefik.http.routers.%s.tls", containerName):              "true",
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", containerName): "le",
	}
}

type CorsSettings struct {
	AllowOrigin string
}

func TraefikCors(settings CorsSettings) map[string]string {
	return map[string]string{
		"traefik.http.middlewares.api-cors.headers.accesscontrolalloworiginlist":  fmt.Sprintf("https://%s", settings.AllowOrigin),
		"traefik.http.middlewares.api-cors.headers.accesscontrolallowmethods":     "GET,POST,PUT,DELETE,OPTIONS",
		"traefik.http.middlewares.api-cors.headers.accesscontrolallowheaders":     "Content-Type,Authorization",
		"traefik.http.middlewares.api-cors.headers.accessControlAllowCredentials": "true",
		"traefik.http.middlewares.api-cors.headers.accesscontrolmaxage":           "100",
		"traefik.http.middlewares.api-cors.headers.addvaryheader":                 "true",
	}
}
