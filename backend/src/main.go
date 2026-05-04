package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if !CheckIfDockerInstalled() {
		log.Fatalf("Docker is not installed or not available in PATH. Please install Docker before running this server.")
		os.Exit(1)
	}

	router := NewRouter()

	var serverPort = os.Getenv("HTTP_PORT")
	if serverPort == "" {
		serverPort = "4000"
	}

	serverPort = FindAvailablePort(serverPort)
	log.Printf("using port %s", serverPort)
	err := http.ListenAndServe(":"+serverPort, router)

	if err != nil {
		log.Fatalf("Server failed: %v. Make sure the port %s is available.", err, serverPort)
		os.Exit(1)
	}
}
