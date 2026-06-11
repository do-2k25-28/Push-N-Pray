package main

import (
	"log"
	"os"

	"pushnpray/cmd/server/api/routes"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/utils"
)

func main() {
	if os.Getenv("REGISTER_TOKEN") == "" {
		log.Fatal("Missing REGISTER_TOKEN environment variable.")
	}

	utils.SetupDockerBad()

	database.InitDB()

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
