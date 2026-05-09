package main

import (
	"log"
	"os"

	"pushnpray/cmd/server/api"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/utils"
)

func main() {
	if !utils.CheckIfDockerInstalled() {
		log.Fatalf("Docker is not installed or not available in PATH. Please install Docker before running this server.")
		os.Exit(1)
	}

	database.InitDB()

	router := api.NewRouter()

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
