package utils

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
)

func ResolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
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

	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer func() {
			_ = l.Close()
		}()
		return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}

	return defaultPort
}
