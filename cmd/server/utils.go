package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
)

func CheckIfDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func FindAvailablePort(defaultPort string) string {
	port, err := strconv.Atoi(defaultPort)
	if err != nil {
		port = 4000
	}

	for port <= 65535 {
		addr := fmt.Sprintf(":%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			_ = l.Close()
			return strconv.Itoa(port)
		}
		port++
	}

	// Fallback to let the OS pick a random available port
	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer func() {
			_ = l.Close()
		}()
		return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}

	return defaultPort
}
