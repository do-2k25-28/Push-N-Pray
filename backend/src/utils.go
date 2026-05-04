package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
)

func CheckIfDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func FindAvailablePort(defaultPort string) string {
	port, err := strconv.Atoi(defaultPort)
	if err != nil {
		port = 4000
	}

	for {
		addr := fmt.Sprintf(":%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			return strconv.Itoa(port)
		}
		port++
	}
}
