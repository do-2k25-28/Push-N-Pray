package main

import (
	"log"
	"os"

	"pushnpray/cmd/server/api/routes"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/utils"
)

func main() {
	database.InitDB()

	cephNetwork, err := utils.InitCephNetwork("pushnpray-ceph")
	if err != nil {
		log.Fatalf("Failed to initialize CEPH network: %v", err)
	}
	log.Printf("CEPH network %q ready (subnet %s, monitor %s)", cephNetwork.Name, cephNetwork.Subnet, cephNetwork.MonitorIP)

	router := routes.NewRouter()

	var serverPort = os.Getenv("HTTP_PORT")
	if serverPort == "" {
		serverPort = "4000"
	}

	serverPort = utils.FindAvailablePort(serverPort)
	log.Printf("using port %s", serverPort)
	err = router.Run(":" + serverPort)

	if err != nil {
		log.Fatalf("Server failed: %v. Make sure the port %s is available.", err, serverPort)
		os.Exit(1)
	}
}
