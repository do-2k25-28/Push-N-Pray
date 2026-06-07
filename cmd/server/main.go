package main

import (
	"context"
	"log"
	"os"

	"pushnpray/cmd/server/api/routes"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/deployment"
	"pushnpray/cmd/server/utils"
)

func main() {
	database.InitDB()

	ctx := context.Background()

	if err := deployment.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize deployment service: %v", err)
	}

	router := routes.NewRouter()

	var serverPort = os.Getenv("HTTP_PORT")
	if serverPort == "" {
		serverPort = "4000"
	}

	serverPort = utils.FindAvailablePort(serverPort)
	log.Printf("using port %s", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatalf("Server failed: %v. Make sure the port %s is available.", err, serverPort)
	}
}
