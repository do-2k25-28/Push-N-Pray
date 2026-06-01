package main

import (
	"context"
	"log"
	"os"

	"pushnpray/cmd/server/api/routes"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/deployment"
	"pushnpray/cmd/server/utils"
	"pushnpray/internal"
)

func main() {
	if !internal.CheckIfDockerInstalled() {
		log.Fatalf("Docker is not installed or not available in PATH. Please install Docker before running this server.")
		os.Exit(1)
	}

	if err := deployment.EnsureTraefik(context.Background()); err != nil {
		log.Fatalf("failed to ensure traefik: %v", err)
		os.Exit(1)
	}

	database.InitDB()

	router := routes.NewRouter()

	var serverPort = os.Getenv("HTTP_PORT")
	if serverPort == "" {
		serverPort = "4000"
	}

	serverPort = utils.FindAvailablePort(serverPort)
	log.Printf("using port %s", serverPort)
	err := router.Run(":" + serverPort)

	if err != nil {
		log.Fatalf("Server failed: %v. Make sure the port %s is available.", err, serverPort)
		os.Exit(1)
	}
}
